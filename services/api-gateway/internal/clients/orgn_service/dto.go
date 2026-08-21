package orgn_service

type GetRequest struct {
	UserID string
	Cursor *string
}

type GetResponse struct {
	Orgns      []*Orgn
	NextCursor *string
}

type Orgn struct {
	ID   string
	Name string
}
