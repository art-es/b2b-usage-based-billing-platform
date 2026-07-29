package auth_service

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	pb "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/grpc/auth_service"
	dto "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/clients/auth_service"
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

func (c *Client) Register(ctx context.Context, dtoReq *dto.RegisterRequest) error {
	req := &pb.RegisterRequest{
		Name:     dtoReq.Name,
		Email:    dtoReq.Email,
		Password: dtoReq.Password,
	}

	_, err := c.client.Register(ctx, req, grpcutil.CallOpts(ctx)...)

	return grpcutil.HandleError(err)
}

func (c *Client) VerifyEmail(ctx context.Context, token string) error {
	req := &pb.VerifyEmailRequest{
		Token: token,
	}

	_, err := c.client.VerifyEmail(ctx, req, grpcutil.CallOpts(ctx)...)

	return grpcutil.HandleError(err)
}

func (c *Client) ResendEmailVerification(ctx context.Context, email string) error {
	req := &pb.ResendEmailVerificationRequest{
		Email: email,
	}

	_, err := c.client.ResendEmailVerification(ctx, req, grpcutil.CallOpts(ctx)...)

	return grpcutil.HandleError(err)
}

func (c *Client) Login(ctx context.Context, dtoReq *dto.LoginRequest) (*dto.LoginResponse, error) {
	req := &pb.LoginRequest{
		Email:    dtoReq.Email,
		Password: dtoReq.Password,
	}

	res, err := c.client.Login(ctx, req, grpcutil.CallOpts(ctx)...)
	if err != nil {
		return nil, grpcutil.HandleError(err)
	}

	return &dto.LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (c *Client) RefreshSession(ctx context.Context, token string) (*dto.RefreshSessionResponse, error) {
	req := &pb.RefreshSessionRequest{
		Token: token,
	}

	res, err := c.client.RefreshSession(ctx, req, grpcutil.CallOpts(ctx)...)
	if err != nil {
		return nil, grpcutil.HandleError(err)
	}

	return &dto.RefreshSessionResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (c *Client) GetMe(ctx context.Context) (*dto.GetMeResponse, error) {
	res, err := c.client.GetMe(ctx, &pb.Empty{}, grpcutil.CallOpts(ctx)...)
	if err != nil {
		return nil, grpcutil.HandleError(err)
	}

	dtoRes := &dto.GetMeResponse{
		SessionID: res.SessionId,
		Name:      res.Name,
		Email:     res.Email,
	}

	if res.Orgn != nil {
		dtoRes.Orgn = &dto.GetMeResponseOrgn{
			ID:   res.Orgn.Id,
			Name: res.Orgn.Name,
		}
	}

	return dtoRes, nil
}
