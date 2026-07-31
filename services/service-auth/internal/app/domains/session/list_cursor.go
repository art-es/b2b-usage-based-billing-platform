package session

import "time"

type ListCursor struct {
	Version   int       `json:"v"`
	ID        string    `json:"i"`
	CreatedAt time.Time `json:"t"`
}

func GetNextCursor(sessions []*Session, hasMore bool) *ListCursor {
	if !hasMore || len(sessions) == 0 {
		return nil
	}

	return &ListCursor{}
}
