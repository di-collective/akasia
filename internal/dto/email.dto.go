package dto

type Email struct {
	From          string    `json:"from"`
	To            string    `json:"to"`
	Subject       string    `json:"subject"`
	TemplateEmail string    `json:"template_email"`
	Body          EmailBody `json:"body"`
}

type EmailBody struct {
	UserName         string `json:"user_name,omitempty"`
	ResetPasswordUrl string `json:"reset_password_url,omitempty"`
	CsMail           string `json:"cs_mail,omitempty"`
}
