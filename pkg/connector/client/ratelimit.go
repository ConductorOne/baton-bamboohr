package client

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const RetryAfterHeader = "Retry-After"

// ExtractRetryAfter - returns a time at which to retry, a bool that is true if retry was found, and if an error occurred.
func ExtractRetryAfter(response *uhttp.WrapperResponse) (*timestamppb.Timestamp, bool, error) {
	retryAfter := response.Header.Get(RetryAfterHeader)
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

func WithBambooHrRatelimitData(resource *v2.RateLimitDescription) uhttp.DoOption {
	return func(response *uhttp.WrapperResponse) error {
		// BambooHR returns a 503 when rate limits are exceeded.
		if response.StatusCode == http.StatusServiceUnavailable {
			resource.Status = v2.RateLimitDescription_STATUS_OVERLIMIT
			// ExtractRetryAfter header used https://documentation.bamboohr.com/docs/api-details
			retryAfter, found, err := ExtractRetryAfter(response)
			if err != nil {
				return err
			}
			if found {
				resource.ResetAt = retryAfter
			}
		}

		return nil
	}
}

func isRatelimited(
	ratelimitStatus v2.RateLimitDescription_Status,
	statusCode int,
) bool {
	return slices.Contains(
		[]v2.RateLimitDescription_Status{
			v2.RateLimitDescription_STATUS_OVERLIMIT,
			v2.RateLimitDescription_STATUS_ERROR,
		},
		ratelimitStatus,
	) || slices.Contains(
		[]int{
			http.StatusTooManyRequests,
			http.StatusGatewayTimeout,
			http.StatusServiceUnavailable,
		},
		statusCode,
	)
}
