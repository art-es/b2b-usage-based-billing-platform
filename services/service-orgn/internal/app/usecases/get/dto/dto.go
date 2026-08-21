package dto

import "github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/domains/orgn"

type Request struct {
	UserID string
	Cursor *string
}

type Response struct {
	List       []*orgn.Orgn
	NextCursor *string
}
