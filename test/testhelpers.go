package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/conductorone/baton-bamboohr/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

func AssertNoRatelimitAnnotations(
	t *testing.T,
	actualAnnotations annotations.Annotations,
) {
	if actualAnnotations != nil && len(actualAnnotations) == 0 {
		return
	}

	for _, annotation := range actualAnnotations {
		var ratelimitDescription v2.RateLimitDescription
		err := annotation.UnmarshalTo(&ratelimitDescription)
		if err != nil {
			continue
		}
		if slices.Contains(
			[]v2.RateLimitDescription_Status{
				v2.RateLimitDescription_STATUS_ERROR,
				v2.RateLimitDescription_STATUS_OVERLIMIT,
			},
			ratelimitDescription.Status,
		) {
			t.Fatal("request was ratelimited, expected not to be ratelimited")
		}
	}
}

func FixturesServer() *httptest.Server {
	return fixturesServerForFile("../../test/fixtures/users_report.json")
}

// FixturesServerWithCustomFields serves a report fixture that includes the
// additional standard fields and a couple of custom fields.
func FixturesServerWithCustomFields() *httptest.Server {
	return fixturesServerForFile("../../test/fixtures/users_report_custom_fields.json")
}

// FixturesServerWithStatus serves the given status code and body for the users
// report request, so tests can exercise error paths (e.g. a 406 for a
// non-existent field).
func FixturesServerWithStatus(statusCode int, body string) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set(uhttp.ContentType, "application/json")
				routeUrl := request.URL.String()
				if !strings.Contains(routeUrl, client.UsersListUrlPath) {
					// This should never happen in tests.
					panic(fmt.Errorf("bad url: %s", routeUrl))
				}
				writer.WriteHeader(statusCode)
				_, _ = writer.Write([]byte(body))
			},
		),
	)
}

func fixturesServerForFile(filename string) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set(uhttp.ContentType, "application/json")
				writer.WriteHeader(http.StatusOK)
				routeUrl := request.URL.String()
				if !strings.Contains(routeUrl, client.UsersListUrlPath) {
					// This should never happen in tests.
					panic(fmt.Errorf("bad url: %s", routeUrl))
				}
				data, _ := os.ReadFile(filename)
				_, err := writer.Write(data)
				if err != nil {
					return
				}
			},
		),
	)
}
