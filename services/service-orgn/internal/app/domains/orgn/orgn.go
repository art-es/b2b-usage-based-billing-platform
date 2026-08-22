package orgn

import "time"

type Orgn struct {
	ID        string
	Name      string
	UserID    string
	CreatedAt time.Time
}

func New(name, userID string) *Orgn {
	return &Orgn{
		Name:   name,
		UserID: userID,
	}
}

func (o *Orgn) Stored() bool {
	return o.ID != ""
}
