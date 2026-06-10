package connector

import (
	"context"
	"testing"

	"github.com/conductorone/baton-bamboohr/pkg/connector/client"
	"github.com/conductorone/baton-bamboohr/test"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/stretchr/testify/require"
)

func TestUsersList(t *testing.T) {
	ctx := context.Background()

	t.Run("should get users with pagination", func(t *testing.T) {
		server := test.FixturesServer()
		defer server.Close()

		bambooClient, err := client.New(
			ctx,
			"mock-access-token",
			"mock-company",
			"",
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		bambooClient.SetBaseUrl(server.URL)
		c := userBuilder(bambooClient)

		resources := make([]*v2.Resource, 0)
		pToken := pagination.Token{
			Token: "",
			Size:  1,
		}
		for {
			nextResources, nextToken, listAnnotations, err := c.List(ctx, nil, &pToken)
			resources = append(resources, nextResources...)

			require.Nil(t, err)
			test.AssertNoRatelimitAnnotations(t, listAnnotations)
			if nextToken == "" {
				break
			}

			pToken.Token = nextToken
		}

		require.NotNil(t, resources)
		require.Len(t, resources, 1)
		require.NotEmpty(t, resources[0].Id)
	})

	t.Run("should include custom fields in user profile", func(t *testing.T) {
		server := test.FixturesServerWithCustomFields()
		defer server.Close()

		bambooClient, err := client.New(
			ctx,
			"mock-access-token",
			"mock-company",
			"",
			[]string{"customJobTitle", "customCostCenter"},
		)
		if err != nil {
			t.Fatal(err)
		}

		bambooClient.SetBaseUrl(server.URL)
		c := userBuilder(bambooClient)

		resources, _, _, err := c.List(ctx, nil, &pagination.Token{})
		require.NoError(t, err)
		require.Len(t, resources, 1)

		userTrait := &v2.UserTrait{}
		annos := resources[0].GetAnnotations()
		for _, a := range annos {
			if a.MessageIs(userTrait) {
				err := a.UnmarshalTo(userTrait)
				require.NoError(t, err)
				break
			}
		}

		profile := userTrait.GetProfile()
		require.NotNil(t, profile)

		fields := profile.GetFields()
		require.NotNil(t, fields)

		jobTitle := fields["customJobTitle"]
		require.NotNil(t, jobTitle)
		require.Equal(t, "Senior Engineer", jobTitle.GetStringValue())

		costCenter := fields["customCostCenter"]
		require.NotNil(t, costCenter)
		require.Equal(t, "ENG-100", costCenter.GetStringValue())
	})
}
