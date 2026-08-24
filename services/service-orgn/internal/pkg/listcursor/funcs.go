package listcursor

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func Encode(secret []byte, obj any) (*string, error) {
	if obj == nil {
		return nil, nil
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("object: json marshal: %w", err)
	}

	hash := generateHash(secret, data)
	out := string(data) + "." + base64.StdEncoding.EncodeToString(hash)
	out = base64.URLEncoding.EncodeToString([]byte(out))

	return &out, nil
}

func DecodeAndCompare(secret []byte, str *string, obj any) error {
	if str == nil {
		return nil
	}

	decoded, err := base64.URLEncoding.DecodeString(*str)
	if err != nil {
		return fmt.Errorf("string: base64 decode: %w", err)
	}

	chunks := strings.Split(string(decoded), ".")
	if len(chunks) != 2 {
		return errors.New("string: must have 2 strings splitted by \".\"")
	}

	data := []byte(chunks[0])
	if err = json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("object: json unmarshal: %w", err)
	}

	hash, err := base64.StdEncoding.DecodeString(chunks[1])
	if err != nil {
		return fmt.Errorf("second chunk: base64 decode: %w", err)
	}

	if !hmac.Equal(hash, generateHash(secret, data)) {
		return errors.New("incorrect hash")
	}

	return nil
}

func generateHash(secret, data []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(data)
	return mac.Sum(nil)
}
