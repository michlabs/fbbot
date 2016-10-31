package fbbot

type User struct {
	ID          string `json:"id"`
	PhoneNumber string `json:"phone_number,omitempty"`
}
