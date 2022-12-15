package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	RetryAfterHeader = "Retry-After"
	APIDomain        = "api.bamboohr.com"
	APIPath          = "api"
	APIGateway       = "gateway.php"
	APIVersion       = "v1"
)

type User struct {
	Id              string `json:"id"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Supervisor      string `json:"supervisor"`
	SupervisorEId   string `json:"supervisorEId"`
	SupervisorId    string `json:"supervisorId"`
	SupervisorEmail string `json:"supervisorEmail"`
	Email           string `json:"workEmail"`
	Status          string `json:"status"`
}

type Fields struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type ReqFields struct {
	Title  string   `json:"title"`
	Fields []string `json:"fields"`
}

type ReportUserResults struct {
	Title  string   `json:"title"`
	Fields []Fields `json:"fields"`
	Users  []*User  `json:"employees"`
}

type BambooHRClient struct {
	Client        *http.Client
	ApiKey        string
	CompanyDomain string
	Transport     *openHttpTransport
}

type UsersResponse struct {
	Users                []*User
	RateLimitDescription *v2.RateLimitDescription
}

type Client interface {
	ListUsers(ctx context.Context, pagination string) (*UsersResponse, error)
}

func New(ctx context.Context, apiKey string, companyDomain string) (*BambooHRClient, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, nil))
	if err != nil {
		return nil, err
	}
	transport := &openHttpTransport{
		base: httpClient.Transport,
	}
	httpClient.Transport = transport
	return &BambooHRClient{
		Client:        httpClient,
		ApiKey:        apiKey,
		Transport:     transport,
		CompanyDomain: companyDomain,
	}, nil
}

type openHttpTransport struct {
	base      http.RoundTripper
	rateLimit *v2.RateLimitDescription
}

func (t *openHttpTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	t.rateLimit = nil // clear previous
	resp, err := t.base.RoundTrip(request)
	if err != nil {
		return resp, err
	}
	// BambooHR returns a 503 when rate limits are exceeded.
	if resp.StatusCode == http.StatusServiceUnavailable {
		t.rateLimit = &v2.RateLimitDescription{
			Status: v2.RateLimitDescription_STATUS_OVERLIMIT,
		}
		// RetryAfter header used https://documentation.bamboohr.com/docs/api-details
		retryAfter, found, err := RetryAfter(resp)
		if err != nil {
			return nil, err
		}
		if found {
			t.rateLimit.ResetAt = retryAfter
		}
	}
	return resp, nil
}

func (c *BambooHRClient) newUnPaginatedURL(path string, v url.Values) (string, error) {
	reqUrl := url.URL{Scheme: "https", Host: APIDomain, Path: strings.Join([]string{APIPath, APIGateway, c.CompanyDomain, APIVersion, path}, "/")}
	reqUrl.RawQuery = v.Encode()
	return reqUrl.String(), nil
}

func (c *BambooHRClient) ListUsers(ctx context.Context) (*UsersResponse, error) {
	users := &ReportUserResults{}
	v := url.Values{}
	v.Set("format", "json")
	reqURL, err := c.newUnPaginatedURL("reports/custom", v)
	if err != nil {
		return nil, err
	}

	listUsersReqBody := ReqFields{
		Title:  "ConductorOne Employees List Report",
		Fields: []string{"firstName", "lastName", "supervisor", "supervisorEId", "supervisorId", "supervisorEmail", "workEmail", "status"},
	}
	body, err := json.Marshal(listUsersReqBody)
	if err != nil {
		return nil, err
	}

	if err := c.query(ctx, http.MethodPost, reqURL, users, body); err != nil {
		return nil, fmt.Errorf("bambooHR-client: error listing users %w", err)
	}
	rv := &UsersResponse{
		Users:                users.Users,
		RateLimitDescription: c.Transport.rateLimit,
	}
	return rv, nil
}

func (c *BambooHRClient) query(ctx context.Context, method string, requestURL string, res interface{}, body []byte) error {
	var bodyReader *bytes.Buffer
	reqUrl, err := url.Parse(requestURL)
	if err != nil {
		return err
	}

	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, reqUrl.String(), bodyReader)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(c.ApiKey, "")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed %d", resp.StatusCode)
	}
	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(rawResp, res); err != nil {
		return err
	}
	return nil
}

// Makes an API call to verify that the given credentials work.
func (c *BambooHRClient) Verify(ctx context.Context) error {
	_, err := c.ListUsers(ctx)
	if err != nil {
		return err
	}

	return nil
}

// RetryAfter - returns a time at which to retry, a bool that is true if retry was found, and if an error occurred.
func RetryAfter(resp *http.Response) (*timestamppb.Timestamp, bool, error) {
	retryAfter := resp.Header.Get(RetryAfterHeader)
	if retryAfter == "" {
		return nil, false, nil
	}
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-retry-after
	// Retry-After = HTTP-date / delay-seconds

	// Retry-After: Fri, 31 Dec 1999 23:59:59 GMT
	if rv, err := time.Parse(time.RFC1123, retryAfter); err == nil {
		return timestamppb.New(rv), true, nil
	}

	// Retry-After: 120
	if duration, err := time.ParseDuration(fmt.Sprintf("%ss", retryAfter)); err == nil {
		return timestamppb.New(time.Now().Add(duration)), true, nil
	}
	return nil, false, fmt.Errorf("unable to parse Retry-After header")
}
