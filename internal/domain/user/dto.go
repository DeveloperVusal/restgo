package user

type UserDto struct {
	Id    int    `json:"id" validate:"omitempty,numeric"`
	Email string `json:"email" validate:"omitempty,email"`
}

type CreateUserDto struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	Name            string `json:"name" validate:"required,alphaunicode"`
	Surname         string `json:"surname" validate:"required,alphaunicode"`
	Activation      bool   `json:"activation" validate:"omitempty,boolean"`
	ConfirmStatus   string `json:"confirm_status" validate:"omitempty,oneof_insensitive=quest waiting success"`
}

type UpdateUserDto struct {
	Id              int    `json:"id" validate:"required,number"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	Name            string `json:"name" validate:"required,alphaunicode"`
	Surname         string `json:"surname" validate:"required,alphaunicode"`
	Activation      bool   `json:"activation" validate:"omitempty,boolean"`
	ConfirmStatus   string `json:"confirm_status" validate:"omitempty,oneof_insensitive=quest waiting success"`
}
