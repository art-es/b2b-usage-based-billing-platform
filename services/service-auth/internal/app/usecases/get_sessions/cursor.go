package get_sessions

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
)

func (u *Usecase) stringToCursor(cursorStr *string) (*session.ListCursor, error) {
	if cursorStr == nil {
		return nil, nil
	}

	splitted := strings.Split(*cursorStr, ".")
	if len(splitted) != 2 {
		return nil, errInvalidCursor
	}

	jsonBytes, err := base64.URLEncoding.DecodeString(splitted[0])
	if err != nil {
		return nil, errInvalidCursor
	}

	signature, err := base64.URLEncoding.DecodeString(splitted[1])
	if err != nil {
		return nil, errInvalidCursor
	}

	var cursor session.ListCursor

	if err = json.Unmarshal(jsonBytes, &cursor); err != nil {
		return nil, errInvalidCursor
	}

	expectedSignature, err := u.keyedHashService.Generate(u.cursorSecretKey, string(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("generate cursor hash: %w", err)
	}

	if string(signature) != expectedSignature {
		return nil, errInvalidCursor
	}

	return &cursor, nil
}

func (u *Usecase) cursorToString(cursor *session.ListCursor) (*string, error) {
	if cursor == nil {
		return nil, nil
	}

	jsonBytes, err := json.Marshal(cursor)
	if err != nil {
		return nil, fmt.Errorf("encode cursor to json: %w", err)
	}

	signature, err := u.keyedHashService.Generate(u.cursorSecretKey, string(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("generate cursor hash: %w", err)
	}

	out := base64.URLEncoding.EncodeToString(jsonBytes) + "." +
		base64.RawStdEncoding.EncodeToString([]byte(signature))

	return &out, nil
}
