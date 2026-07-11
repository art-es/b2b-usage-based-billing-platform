package broker

const (
	SubjectEmailSend = "email.send"
)

type ProduceMessage struct {
	Subject        string
	IdempotencyKey string
	Payload        []byte
}

func SupportedSubjects() []string {
	return []string{SubjectEmailSend}
}
