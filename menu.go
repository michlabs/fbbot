package fbbot

type Menu struct {
	Locale                string     `json:"locale"`
	ComposerInputDisabled bool       `json:"composer_input_disabled"`
	CallToActions         []*MenuItem `json:"call_to_actions"`
}

func NewMenu() *Menu {
	return &Menu{
		Locale: "default",
	}
}

func (m *Menu) AddMenuItems(items ...*MenuItem) {
	m.CallToActions = append(m.CallToActions, items...)
}

type MenuItem struct {
	Title              string `json:"title"`
	Type               string `json:"type"`
	URL                string `json:"url,omitempty"`
	WebviewHeightRatio string `json:"webview_height_ratio,omitempty"`
	Payload string `json:"payload,omitempty"`
	CallToActions []*MenuItem `json:"call_to_actions,omitempty"`
}

func NewWebURLMenuItem(title, url string) *MenuItem {
	return &MenuItem{
		Title:              title,
		Type:               "web_url",
		URL:                url,
		WebviewHeightRatio: "full",
	}
}

func NewPostbackMenuItem(title, payload string) *MenuItem {
	return &MenuItem{
		Title:   title,
		Type:    "postback",
		Payload: payload,
	}
}

func NewNestedMenuItem(title string) *MenuItem {
	return &MenuItem{
		Title: title,
		Type:  "nested",
	}
}

// AddMenuItems only used for nested menu item
func (mi *MenuItem) AddMenuItems(items ...*MenuItem) {
	if mi.Type == "nested" {
		mi.CallToActions = append(mi.CallToActions, items...)
	}
}
