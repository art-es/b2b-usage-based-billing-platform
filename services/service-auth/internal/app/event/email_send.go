package event

type EmailSend struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}
