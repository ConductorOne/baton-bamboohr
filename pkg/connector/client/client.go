package client

import (
	"context"
	"errors"
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
	ListUsers(ctx context.Context) ([]*User, *v2.RateLimitDescription, error)
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

	fields := make([]string, 0, len(StandardFields)+len(c.CustomFields))
	fields = append(fields, StandardFields...)
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
		// BambooHR rejects the whole report with 406 Not Acceptable when the
		// request references a field that doesn't exist in the schema. The
		// standard fields always exist, so in practice a 406 means one of the
		// configured custom-fields is misspelled or isn't a real BambooHR field
		// alias. validateCustomFields can't catch this — it only runs on a
		// successful (200) response — so surface a clear, actionable error here.
		var reqErr *RequestError
		if len(c.CustomFields) > 0 && errors.As(err, &reqErr) && reqErr.Status == http.StatusNotAcceptable {
			return nil, ratelimitData, fmt.Errorf(
				"bambooHR-client: report rejected (406 Not Acceptable) for a non-existent field. "+
					"One of the configured custom-fields %v is not a valid BambooHR field alias; "+
					"check the field names against the BambooHR field-name reference: %w",
				c.CustomFields, err,
			)
		}
		return nil, ratelimitData, fmt.Errorf("bambooHR-client: error listing users %w", err)
	}

	if err := c.validateCustomFields(users.Fields); err != nil {
		return nil, ratelimitData, err
	}

	return users.Users, ratelimitData, nil
}

// validateCustomFields fails fast if any configured custom field was not
// returned by the report, so misconfigured field names surface as a clear
// error instead of silently-empty profile attributes.
func (c *BambooHRClient) validateCustomFields(returned []Fields) error {
	if len(c.CustomFields) == 0 {
		return nil
	}

	present := make(map[string]bool, len(returned))
	for _, f := range returned {
		present[f.Id] = true
	}

	var missing []string
	for _, configured := range c.CustomFields {
		if !present[configured] {
			missing = append(missing, configured)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf(
			"bambooHR-client: configured custom-fields not found in report response: %v. "+
				"Check the field names against the BambooHR field-name reference",
			missing,
		)
	}
	return nil
}

// Verify - Makes an API call to verify that the given credentials work.
func (c *BambooHRClient) Verify(ctx context.Context) error {
	_, _, err := c.ListUsers(ctx)
	return err
}
