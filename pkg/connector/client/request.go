package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	APIDomain                 = "api.bamboohr.com"
	APIPath                   = "api"
	APIGateway                = "gateway.php"
	APIVersion                = "v1"
	BambooPasswordPlaceholder = "x"
)

type RequestError struct {
	Status int
	URL    *url.URL
	Body   string
}

func (r *RequestError) Error() string {
	return fmt.Sprintf(
		"bamboohr-connector: request error. Status: %d, Url: %s, Body: %s",
		r.Status,
		r.URL,
		r.Body,
	)
}

func (c *BambooHRClient) newUnPaginatedURL(path string, v url.Values) *url.URL {
	return &url.URL{
		Scheme: c.BaseUrl.Scheme,
		Host:   c.BaseUrl.Host,
		Path: strings.Join(
			[]string{
				APIPath,
				APIGateway,
				c.CompanyDomain,
				APIVersion,
				path,
			},
			"/",
		),
		RawQuery: v.Encode(),
	}
}

func (c *BambooHRClient) makeRequest(
	ctx context.Context,
	url *url.URL,
	target interface{},
	method string,
	requestBody io.Reader,
) (*v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, method, url.String(), requestBody)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.ApiKey, BambooPasswordPlaceholder)

	ratelimitData := v2.RateLimitDescription{}

	response, err := c.wrapper.Do(
		req,
		WithBambooHrRatelimitData(&ratelimitData),
		uhttp.WithJSONResponse(target),
	)
	if err == nil {
		return &ratelimitData, nil
	}
	if response == nil {
		return nil, err
	}
	defer response.Body.Close()

	// If we get ratelimit data back (e.g. the "Retry-After" header) or a
	// "ratelimit-like" status code, then return a recoverable gRPC code.
	if isRatelimited(ratelimitData.Status, response.StatusCode) {
		return &ratelimitData, status.Error(codes.Unavailable, response.Status)
	}

	// If it's some other error, it is unrecoverable.
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return nil, &RequestError{
		URL:    url,
		Status: response.StatusCode,
		Body:   string(responseBody),
	}
}
