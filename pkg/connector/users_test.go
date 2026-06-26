package connector

import (
	"context"
	"testing"

	"github.com/conductorone/baton-bamboohr/pkg/connector/client"
	"github.com/conductorone/baton-bamboohr/test"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
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

	t.Run("should include the additional standard fields in user profile", func(t *testing.T) {
		server := test.FixturesServerWithCustomFields()
		defer server.Close()

		bambooClient, err := client.New(ctx, "mock-access-token", "mock-company", "", nil)
		require.NoError(t, err)
		bambooClient.SetBaseUrl(server.URL)

		resources, _, _, err := userBuilder(bambooClient).List(ctx, nil, &pagination.Token{})
		require.NoError(t, err)
		require.Len(t, resources, 1)

		profile := profileFields(t, resources[0])
		require.Equal(t, "Jane", profile["firstName"].GetStringValue())
		require.Equal(t, "Doe", profile["lastName"].GetStringValue())
		require.Equal(t, "Senior Engineer", profile["jobTitle"].GetStringValue())
		require.Equal(t, "2020-01-15", profile["hireDate"].GetStringValue())
		require.Equal(t, "E-101", profile["employeeNumber"].GetStringValue())
		require.Equal(t, "Remote - US", profile["location"].GetStringValue())
		require.Equal(t, "Janie", profile["preferredName"].GetStringValue())
		require.Equal(t, "Jane Doe", profile["displayName"].GetStringValue())
		require.Equal(t, "2018-03-01", profile["originalHireDate"].GetStringValue())
		require.Equal(t, "John Manager", profile["reportsTo"].GetStringValue())
	})

	t.Run("should include configured custom fields in user profile", func(t *testing.T) {
		server := test.FixturesServerWithCustomFields()
		defer server.Close()

		bambooClient, err := client.New(
			ctx,
			"mock-access-token",
			"mock-company",
			"",
			[]string{"customField4444", "customCostCenter"},
		)
		require.NoError(t, err)
		bambooClient.SetBaseUrl(server.URL)

		resources, _, _, err := userBuilder(bambooClient).List(ctx, nil, &pagination.Token{})
		require.NoError(t, err)
		require.Len(t, resources, 1)

		profile := profileFields(t, resources[0])
		require.Equal(t, "Ace", profile["customField4444"].GetStringValue())
		require.Equal(t, "ENG-100", profile["customCostCenter"].GetStringValue())
	})

	t.Run("should error when a configured custom field is missing from the response", func(t *testing.T) {
		server := test.FixturesServerWithCustomFields()
		defer server.Close()

		bambooClient, err := client.New(
			ctx,
			"mock-access-token",
			"mock-company",
			"",
			[]string{"doesNotExist"},
		)
		require.NoError(t, err)
		bambooClient.SetBaseUrl(server.URL)

		_, _, _, err = userBuilder(bambooClient).List(ctx, nil, &pagination.Token{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "doesNotExist")
	})
}

// profileFields pulls the user-trait profile field map off a synced resource.
func profileFields(t *testing.T, resource *v2.Resource) map[string]*structpb.Value {
	t.Helper()

	userTrait := &v2.UserTrait{}
	for _, a := range resource.GetAnnotations() {
		if a.MessageIs(userTrait) {
			require.NoError(t, a.UnmarshalTo(userTrait))
			break
		}
	}
	profile := userTrait.GetProfile()
	require.NotNil(t, profile)
	return profile.GetFields()
}
