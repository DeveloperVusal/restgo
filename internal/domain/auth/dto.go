package auth

type UserDto struct {
	Id    int    `json:"id" validate:"omitempty,numeric"`
	Email string `json:"email" validate:"omitempty,email"`
}

type LoginDto struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Device    string
	Ip        string
	UserAgent string
}

type LogoutDto struct {
	Id    int
	Token string
}

type RegistrationDto struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	Name            string `json:"name" validate:"required,alphaunicode"`
	Surname         string `json:"surname" validate:"required,alphaunicode"`
}
