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

type CreateRequest struct {
	UserID string
	Name   string
}

type CreateResponse struct {
	OrgnID string
}

type UpdateRequest struct {
	UserID string
	OrgnID string
	Name   *string
}
