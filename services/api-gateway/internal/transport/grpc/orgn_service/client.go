package orgn_service

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	dto "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/clients/orgn_service"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/grpc/orgn_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/grpc/grpcutil"
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

func (c *Client) Get(ctx context.Context, dtoReq *dto.GetRequest) (*dto.GetResponse, error) {
	req := &pb.GetRequest{
		UserId: dtoReq.UserID,
		Cursor: dtoReq.Cursor,
	}

	res, err := c.client.Get(ctx, req, grpcutil.CallOpts(ctx)...)
	if err != nil {
		return nil, grpcutil.HandleError(err)
	}

	dtoOrgns := make([]*dto.Orgn, 0, len(res.Orgns))
	for _, orgn := range res.Orgns {
		dtoOrgns = append(dtoOrgns, &dto.Orgn{
			ID:   orgn.Id,
			Name: orgn.Name,
		})
	}

	return &dto.GetResponse{
		Orgns:      dtoOrgns,
		NextCursor: res.NextCursor,
	}, nil
}
