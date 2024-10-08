package dns

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidType   = errors.New("invalid type")
	ErrServerError   = errors.New("server error")
	ErrInvalidToken  = errors.New("invalid token")
	ErrAlreadyExists = errors.New("item already exists")
)

// Record is a representation of DNS record.
type Record struct {
	ID     string
	Name   string
	TTL    time.Duration
	Type   string
	Values []string
}

type Client interface {
	GetRecords(ctx context.Context, domain string) ([]*Record, error)
	CreateRecord(
		ctx context.Context,
		domain, name, recordType string,
		values []string,
		ttl time.Duration,
	) error
	ReplaceRecord(
		ctx context.Context,
		domain, name, recordType string,
		values []string,
		ttl time.Duration,
	) error
}
