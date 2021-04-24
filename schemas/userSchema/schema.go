package userSchema

type UserSchema struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Username    string `json:"username"`
	DateCreated string `json:"date_created"`
	DarkMode    bool   `json:"dark_mode"`
}
