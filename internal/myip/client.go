package myip

import (
	"context"
)

type Client interface {
	GetIPAddress(ctx context.Context) (string, error)
}
