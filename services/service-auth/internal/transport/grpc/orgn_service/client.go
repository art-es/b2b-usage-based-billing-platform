package orgn_service

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/orgn"
	dto "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/clients/orgn_service"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/generated/grpc/orgn_service"
)

type Client struct {
	io.Closer
	client pb.OrgnServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return &Client{
		Closer: conn,
		client: pb.NewOrgnServiceClient(conn),
	}, nil
}

func (c *Client) GetByID(ctx context.Context, req *dto.GetByIDRequest) (*orgn.Orgn, error) {
	org, err := c.client.GetById(ctx, &pb.GetByIdRequest{
		UserId: req.UserID,
		OrgnId: req.OrgnID,
	})
	if err != nil {
		return nil, err
	}

	return &orgn.Orgn{
		ID:     org.Id,
		Name:   org.Name,
		UserID: org.UserId,
	}, nil
}
