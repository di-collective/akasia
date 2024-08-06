package repository

type tables struct {
	User          string
	Profile       string
	Message       string
	UserMessage   string
	ResetPassword string
	Clinic        string
	Location      string
	Event         string
}

type views struct {
	UserMessage string
}

var (
	Tables = tables{
		User:          "user",
		Profile:       "profile",
		Message:       "message",
		UserMessage:   "user_message",
		ResetPassword: "reset_password",
		Clinic:        "clinic",
		Location:      "location",
		Event:         "event",
	}
	Views = views{
		UserMessage: "view_user_message",
	}
)
