package orgn

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	DBListLimit       = ResponseListLimit + 1
	ResponseListLimit = 10
)

type ListCursor struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func HandleList(list []*Orgn) ([]*Orgn, *ListCursor) {
	if len(list) < DBListLimit {
		return list, nil
	}

	last := list[len(list)-1]
	list = list[:len(list)-1]

	return list, &ListCursor{
		ID:        last.ID,
		CreatedAt: last.CreatedAt,
	}
}

func NewListCursorFromString(secret []byte, str *string) (*ListCursor, error) {
	if str == nil {
		return nil, nil
	}

	decoded, err := base64.URLEncoding.DecodeString(*str)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 format: %w", err)
	}

	chunks := strings.Split(string(decoded), ".")
	if len(chunks) != 2 {
		return nil, errors.New("invalid cursor format: must have 2 chunks")
	}

	data := []byte(chunks[0])

	var obj ListCursor
	if err = json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("invalid json cursor: %w", err)
	}

	hash, err := base64.StdEncoding.DecodeString(chunks[1])
	if err != nil {
		return nil, fmt.Errorf("invalid base64 format of hash: %w", err)
	}

	if !hmac.Equal(hash, generateHash(secret, data)) {
		return nil, errors.New("incorrect hash")
	}

	return &obj, nil
}

func (c *ListCursor) String(secret []byte) (*string, error) {
	if c == nil {
		return nil, nil
	}

	data, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("json marshal cursor: %w", err)
	}

	hash := generateHash(secret, data)
	out := string(data) + "." + base64.StdEncoding.EncodeToString(hash)
	out = base64.URLEncoding.EncodeToString([]byte(out))

	return &out, nil
}

func generateHash(secret, data []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(data)
	return mac.Sum(nil)
}
