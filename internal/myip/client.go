package myip

import (
	"context"
	"errors"
)

var ErrInvalidType = errors.New("invalid type")

type Client interface {
	GetIPAddress(ctx context.Context) (string, error)
}
