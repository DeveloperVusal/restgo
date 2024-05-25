package user

type UserDto struct {
	Id    int    `json:"id" validate:"omitempty,numeric"`
	Email string `json:"email" validate:"omitempty,email"`
}
