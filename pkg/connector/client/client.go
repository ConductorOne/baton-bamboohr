package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
}

type Client interface {
	ListUsers(ctx context.Context, pagination string) ([]*User, *v2.RateLimitDescription, error)
}

func New(ctx context.Context, apiKey string, companyDomain string) (*BambooHRClient, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, nil))
	if err != nil {
		return nil, err
	}
	wrapper := uhttp.NewBaseHttpClient(httpClient)

	baseUrl := url.URL{
		Scheme: "https",
		Host:   APIDomain,
	}
	return &BambooHRClient{
		wrapper:       wrapper,
		ApiKey:        apiKey,
		CompanyDomain: companyDomain,
		BaseUrl:       &baseUrl,
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
	reqURL := c.newUnPaginatedURL(UsersListUrlPath, v)

	listUsersReqBody := ReqFields{
		Title: "ConductorOne Employees List Report",
		Fields: []string{
			"firstName",
			"lastName",
			"supervisor",
			"supervisorEId",
			"supervisorId",
			"supervisorEmail",
			"workEmail",
			"status",
		},
	}
	bodyBytes, err := json.Marshal(listUsersReqBody)
	if err != nil {
		return nil, nil, err
	}
	body := strings.NewReader(string(bodyBytes))

	ratelimitData, err := c.makeRequest(
		ctx,
		reqURL,
		users,
		http.MethodPost,
		body,
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
