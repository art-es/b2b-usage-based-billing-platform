package listcursor

type Encoder[T any] struct {
	secret []byte
}

func NewEncoder[T any](secret []byte) *Encoder[T] {
	return &Encoder[T]{
		secret: secret,
	}
}

func (e *Encoder[T]) Encode(obj *T) (*string, error) {
	return Encode(e.secret, obj)
}

func (e *Encoder[T]) DecodeAndCompare(str *string) (*T, error) {
	var obj T

	err := DecodeAndCompare(e.secret, str, &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
