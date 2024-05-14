package dto

type RequestRegisterUser struct {
	Provider       string `json:"provider" validate:"required,oneof=email google.com facebook.com apple.com"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required_if=Provider email,min=12,max=128"`
	RepeatPassword string `json:"repeat_password" validate:"required,eqfield=Password"`
}
