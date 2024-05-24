package dto

type RequestRegisterUser struct {
	Provider       string `json:"provider" validate:"required,oneof=email google.com facebook.com apple.com"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required_if=Provider email,min=12,max=128"`
	RepeatPassword string `json:"repeat_password" validate:"required,eqfield=Password"`
}

type RequestForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type RequestUpdatePassword struct {
	UserID     string `json:"user_id" validate:"required"`
	ResetToken string `json:"reset_token" validate:"required"`
	Password   string `json:"password" validate:"required_if=Provider email,min=12,max=128"`
}
