package fbbot

type User struct {
	ID               string `json:"id"`
	PhoneNumber      string `json:"phone_number,omitempty"`
	isFetched        bool
	firstName        string `json:"first_name, omitempty"`
	lastName         string `json:"last_name, omitempty"`
	profilePic       string `json:"profile_pic, omitempty"`
	locale           string `json:"locale, omitempty"`
	timezone         int    `json:"timezone, omitempty"`
	gender           string `json:"gender, omitempty"`
	isPaymentEnabled bool   `json:"is_payment_enabled, omitempty"` // Is the user eligible to receive messenger platform payment messages
}

func (u *User) FirstName() string {
	if !u.isFetched {
		bot.fetchUserData(u)
	}
	return u.firstName
}

func (u *User) LastName() string {
	if !u.isFetched {
		bot.fetchUserData(u)
	}
	return u.lastName
}

func (u *User) FullName() string {
	return u.FirstName() + " " + u.LastName()
}

func (u *User) ProfilePic() string {
	if !u.isFetched {
		bot.fetchUserData(u)
	}
	return u.profilePic
}

func (u *User) Locale() string {
	if !u.isFetched {
		bot.fetchUserData(u)
	}
	return u.locale
}

func (u *User) Timezone() int {
	if !u.isFetched {
		bot.fetchUserData(u)
	}
	return u.timezone
}

func (u *User) Gender() string {
	if !u.isFetched {
		bot.fetchUserData(u)
	}

	return u.gender
}

func (u *User) IsPaymentEnabled() bool {
	if !u.isFetched {
		bot.fetchUserData(u)
	}
	return u.isPaymentEnabled
}
