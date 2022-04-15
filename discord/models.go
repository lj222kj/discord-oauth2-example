package discord

type User struct {
	UserId   string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
