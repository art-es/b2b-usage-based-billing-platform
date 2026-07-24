package auth_service

import "context"

type Client struct{}

func NewClient(addr string) *Client {
	return &Client{}
}

func (c *Client) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	return nil, nil
}
