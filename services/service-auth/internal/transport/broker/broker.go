package broker

type ProduceMessage struct {
	Topic string
	Key   []byte
	Value []byte
}
