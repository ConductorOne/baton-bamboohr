package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

const (
	UsersListUrlPath = "reports/custom"
)

type BambooHRClient struct {
	wrapper       *uhttp.BaseHttpClient
	ApiKey        string
	CompanyDomain string
	BaseUrl       *url.URL
	CustomFields  []string
}

type Client interface {
	ListUsers(ctx context.Context, pagination string) ([]*User, *v2.RateLimitDescription, error)
}

func New(ctx context.Context, apiKey string, companyDomain string, baseURL string, customFields []string) (*BambooHRClient, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, nil))
	if err != nil {
		return nil, err
	}
	wrapper := uhttp.NewBaseHttpClient(httpClient)

	var parsedBaseURL *url.URL
	if baseURL != "" {
		parsedBaseURL, err = url.Parse(baseURL)
		if err != nil {
			return nil, fmt.Errorf("invalid base URL: %w", err)
		}
		if parsedBaseURL.Scheme != "http" && parsedBaseURL.Scheme != "https" {
			return nil, fmt.Errorf("base-url must have http or https scheme, got %q", parsedBaseURL.Scheme)
		}
	} else {
		parsedBaseURL = &url.URL{
			Scheme: "https",
			Host:   APIDomain,
		}
	}
	return &BambooHRClient{
		wrapper:       wrapper,
		ApiKey:        apiKey,
		CompanyDomain: companyDomain,
		BaseUrl:       parsedBaseURL,
		CustomFields:  customFields,
	}, nil
}

// SetBaseUrl shim for local integration tests.
func (c *BambooHRClient) SetBaseUrl(rawUrl string) {
	baseUrl, err := url.Parse(rawUrl)
	if err != nil {
		return
	}
	c.BaseUrl = baseUrl
}

func (c *BambooHRClient) ListUsers(ctx context.Context) (
	[]*User,
	*v2.RateLimitDescription,
	error,
) {
	users := &ReportUserResults{}
	v := url.Values{}
	v.Set("format", "json")
	v.Set("onlyCurrent", "false")
	reqURL := c.newUnPaginatedURL(UsersListUrlPath, v)

	fields := []string{
		"firstName",
		"lastName",
		"supervisor",
		"supervisorEId",
		"supervisorId",
		"supervisorEmail",
		"workEmail",
		"status",
		"department",
		"division",
		"employmentHistoryStatus",
		"terminationDate",
	}
	fields = append(fields, c.CustomFields...)

	listUsersReqBody := ReqFields{
		Title:  "ConductorOne Employees List Report",
		Fields: fields,
	}

	ratelimitData, err := c.makeRequest(
		ctx,
		reqURL,
		users,
		http.MethodPost,
		listUsersReqBody,
	)
	if err != nil {
		return nil, ratelimitData, fmt.Errorf("bambooHR-client: error listing users %w", err)
	}
	return users.Users, ratelimitData, nil
}

// Verify - Makes an API call to verify that the given credentials work.
func (c *BambooHRClient) Verify(ctx context.Context) error {
	_, _, err := c.ListUsers(ctx)
	return err
}
