package rmclient

import (
	"context"
	"errors"
	"io"
)

type Client struct{}

type AuthResponse struct {
	Token  string
	UserID string
}

func (c *Client) Auth(ctx context.Context, code string) (*AuthResponse, error) {
	return nil, errors.New("implement me")
}

func (c *Client) Upload(ctx context.Context, name string, r io.Reader) error {
	return errors.New("not implemented")
}
