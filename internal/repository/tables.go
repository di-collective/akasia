package repository

type tables struct {
	User        string
	Profile     string
	Message     string
	UserMessage string
}

type views struct {
	UserMessage string
}

var (
	Tables = tables{
		User:        "user",
		Profile:     "string",
		Message:     "message",
		UserMessage: "user_message",
	}
	Views = views{
		UserMessage: "view_user_message",
	}
)
