package session

import (
	"time"
)

const (
	DBListLimit       = ResponseListLimit + 1
	ResponseListLimit = 10
)

type ListCursor struct {
	Version   int       `json:"v"`
	ID        string    `json:"i"`
	CreatedAt time.Time `json:"t"`
}

func HandleList(list []*Session) ([]*Session, *ListCursor) {
	if len(list) < DBListLimit {
		return list, nil
	}

	last := list[len(list)-1]
	list = list[:len(list)-1]

	return list, &ListCursor{
		Version:   1,
		ID:        last.ID,
		CreatedAt: last.CreatedAt,
	}
}
