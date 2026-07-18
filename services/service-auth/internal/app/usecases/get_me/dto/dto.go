package dto

type Response struct {
	SessionID string
	User      ResponseUser
	Orgn      *ResponseOrgn
}

type ResponseUser struct {
	Email string
	Name  string
}

type ResponseOrgn struct {
	ID   string
	Name string
}
