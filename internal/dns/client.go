package dns

import (
	"context"
	"fmt"
	"time"
)

// Record is a representation of DNS record.
type Record struct {
	ID     string
	Name   string
	TTL    time.Duration
	Type   string
	Values []string
}

var ErrAlreadyExists = fmt.Errorf("item already exists")

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
