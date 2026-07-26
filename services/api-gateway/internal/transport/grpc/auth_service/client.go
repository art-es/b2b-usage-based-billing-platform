package auth_service

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	dto "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/clientdto/auth_service"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/grpc/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/grpc/grpcutil"
)

type Client struct {
	io.Closer
	client pb.AuthServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return &Client{
		Closer: conn,
		client: pb.NewAuthServiceClient(conn),
	}, nil
}

func (c *Client) Register(ctx context.Context, req *dto.RegisterRequest) error {
	_, err := c.client.Register(ctx, &pb.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	return grpcutil.HandleError(err)
}

func (c *Client) VerifyEmail(ctx context.Context, token string) error {
	_, err := c.client.VerifyEmail(ctx, &pb.VerifyEmailRequest{
		Token: token,
	})
	return grpcutil.HandleError(err)
}

func (c *Client) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	res, err := c.client.Login(ctx, &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, grpcutil.HandleError(err)
	}

	return &dto.LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}
