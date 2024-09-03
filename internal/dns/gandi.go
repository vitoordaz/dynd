package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
)

const (
	gandiTokenInfoEndpoint = "https://id.gandi.net/tokeninfo"
	gandiRecordsEndpoint   = "https://api.gandi.net/v5/livedns/domains/%s/records"
	gandiRecordEndpoint    = "https://api.gandi.net/v5/livedns/domains/%s/records/%s/%s"
)

func NewGandiClient(ctx context.Context, accessToken string) (Client, error) {
	client := resty.New().
		SetRetryCount(5).
		SetAuthScheme("Bearer").
		SetAuthToken(accessToken).
		SetJSONUnmarshaler(json.Unmarshal)
	if err := validateAccessToken(ctx, client); err != nil {
		return nil, err
	}
	return &gandiClient{client}, nil
}

func getRequiredScope() []string {
	return []string{"domain:view", "domain:tech"}
}

type gandiClient struct {
	client *resty.Client
}

type gandiRecord struct {
	HREF   string   `json:"rrset_href,omitempty"`
	Name   string   `json:"rrset_name,omitempty"`
	TTL    int64    `json:"rrset_ttl,omitempty"`
	Type   string   `json:"rrset_type,omitempty"`
	Values []string `json:"rrset_values,omitempty"`
}

type gandiError struct {
	Cause   string `json:"cause"`
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Object  string `json:"object"`
}

func (c *gandiClient) GetRecords(ctx context.Context, domain string) ([]*Record, error) {
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult([]*gandiRecord{}).
		Get(fmt.Sprintf(gandiRecordsEndpoint, url.QueryEscape(domain)))
	if err != nil {
		return nil, err
	}
	gandiRecords, ok := resp.Result().(*[]*gandiRecord)
	if !ok {
		return nil, fmt.Errorf("invalid result type: %s", reflect.TypeOf(resp.Result()))
	}
	records := make([]*Record, 0, len(*gandiRecords))
	for _, r := range *gandiRecords {
		records = append(records, &Record{
			ID:     r.HREF,
			Name:   r.Name,
			TTL:    time.Duration(r.TTL),
			Type:   r.Type,
			Values: r.Values,
		})
	}
	return records, nil
}

func (c *gandiClient) CreateRecord(
	ctx context.Context,
	domain, name, recordType string,
	values []string,
	ttl time.Duration,
) error {
	record := &gandiRecord{
		TTL:    int64(ttl / time.Second),
		Values: values,
	}
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(record).
		SetResult(&gandiError{}).
		Post(fmt.Sprintf(
			gandiRecordEndpoint,
			url.QueryEscape(domain),
			url.QueryEscape(name),
			url.QueryEscape(recordType),
		))
	if err != nil {
		return err
	}
	if status := resp.StatusCode(); status == http.StatusCreated {
		return nil
	} else if status == http.StatusOK {
		return ErrAlreadyExists
	}
	result, ok := resp.Result().(*gandiError)
	if !ok {
		return fmt.Errorf("invalid result type: %s", reflect.TypeOf(resp.Result()))
	}
	return fmt.Errorf(result.Message)
}

func (c *gandiClient) ReplaceRecord(
	ctx context.Context,
	domain, name, recordType string,
	values []string,
	ttl time.Duration,
) error {
	record := &gandiRecord{
		TTL:    int64(ttl / time.Second),
		Values: values,
	}
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(record).
		SetResult(&gandiError{}).
		Put(fmt.Sprintf(
			gandiRecordEndpoint,
			url.QueryEscape(domain),
			url.QueryEscape(name),
			url.QueryEscape(recordType),
		))
	if err != nil {
		return err
	}
	status := resp.StatusCode()
	if status == http.StatusCreated {
		return nil
	}
	if status == http.StatusOK {
		return ErrAlreadyExists
	}
	result, ok := resp.Result().(*gandiError)
	if !ok {
		return fmt.Errorf("invalid result type: %s", reflect.TypeOf(resp.Result()))
	}
	return fmt.Errorf(result.Message)
}

type tokenInfoResponse struct {
	Scope []string `json:"scope"`
}

// validateAccessToken validates that client's access token has required permissions.
func validateAccessToken(ctx context.Context, client *resty.Client) error {
	resp, err := client.R().
		SetContext(ctx).
		SetResult(&tokenInfoResponse{}).
		Get(gandiTokenInfoEndpoint)
	if err != nil {
		return err
	}
	result, ok := resp.Result().(*tokenInfoResponse)
	if !ok {
		return fmt.Errorf("invalid result type: %s", reflect.TypeOf(result))
	}
	if err := validateScope(result.Scope, getRequiredScope()); err != nil {
		return err
	}
	return nil
}

// validateScope validates that a given scope from a token includes required scope.
func validateScope(tokenScope []string, required []string) error {
	grantedScope := make(map[string]bool, len(tokenScope))
	for _, scope := range tokenScope {
		grantedScope[scope] = true
	}
	var missingScope []string
	for _, scope := range required {
		if !grantedScope[scope] {
			missingScope = append(missingScope, scope)
		}
	}
	if len(missingScope) > 0 {
		return fmt.Errorf("token is missing following scope: %s", strings.Join(missingScope, ", "))
	}
	return nil
}
