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

type DestroyDto struct {
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

type RefreshDto struct {
	Refresh   string
	Device    string
	Ip        string
	UserAgent string
}

type ActivationDto struct {
	Email string `json:"email" validate:"required,email"`
	Key   string `json:"key" validate:"required,sha256"`
	Code  int    `json:"code" validate:"required,numeric"`
}
