package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-bamboohr/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
)

type UserResourceType struct {
	resourceType   *v2.ResourceType
	bambooHRClient *client.BambooHRClient
}

func (o *UserResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return o.resourceType
}

func (o *UserResourceType) List(
	ctx context.Context,
	_ *v2.ResourceId,
	pt *pagination.Token,
) ([]*v2.Resource, string, annotations.Annotations, error) {
	usersResponse, err := o.bambooHRClient.ListUsers(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	rv := make([]*v2.Resource, 0, len(usersResponse.Users))
	for _, user := range usersResponse.Users {
		newResource, err := userResource(ctx, user)
		if err != nil {
			return nil, "", nil, err
		}
		rv = append(rv, newResource)
	}

	return rv, "", nil, nil
}

func (o *UserResourceType) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *UserResourceType) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func userBuilder(bambooHRClient *client.BambooHRClient) *UserResourceType {
	return &UserResourceType{
		resourceType:   resourceTypeUser,
		bambooHRClient: bambooHRClient,
	}
}

// userResource convert a BambooHR into a Resource.
func userResource(
	ctx context.Context,
	user *client.User,
) (*v2.Resource, error) {
	profile := userProfile(ctx, user)
	displayName := fmt.Sprintf(
		"%s %s",
		user.FirstName,
		user.LastName,
	)
	userTraitOptions := []resource.UserTraitOption{
		resource.WithUserProfile(profile),
		resource.WithEmail(user.Email, true),
	}

	return resource.NewUserResource(
		displayName,
		resourceTypeUser,
		user.Id,
		userTraitOptions,
	)
}

func userProfile(ctx context.Context, user *client.User) map[string]interface{} {
	profile := make(map[string]interface{})
	profile["supervisorEId"] = user.SupervisorEId
	profile["supervisorFullName"] = user.Supervisor
	profile["supervisorId"] = user.SupervisorId
	profile["supervisorEmail"] = user.SupervisorEmail
	profile["user_id"] = user.Id

	return profile
}
