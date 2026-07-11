package event

type EmailSend struct {
	IdempotencyKey string
	Email          string
	Subject        string
	Content        string
}
