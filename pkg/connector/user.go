package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-bamboohr/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/sdk"
)

type UserResourceType struct {
	resourceType   *v2.ResourceType
	bambooHRClient *client.BambooHRClient
}

func (o *UserResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return o.resourceType
}

func (o *UserResourceType) List(ctx context.Context, _ *v2.ResourceId, pt *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	usersResponse, err := o.bambooHRClient.ListUsers(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	rv := make([]*v2.Resource, 0, len(usersResponse.Users))
	for _, user := range usersResponse.Users {
		annos := &v2.V1Identifier{
			Id: user.Id,
		}
		profile := userProfile(ctx, user)
		userResource, err := sdk.NewUserResource(fmt.Sprintf("%s %s", user.FirstName, user.LastName), resourceTypeUser, nil, user.Id, user.Email, profile, annos)
		if err != nil {
			return nil, "", nil, err
		}
		rv = append(rv, userResource)
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

func userProfile(ctx context.Context, user *client.User) map[string]interface{} {
	profile := make(map[string]interface{})
	profile["supervisorEId"] = user.SupervisorEId
	profile["supervisorFullName"] = user.Supervisor
	profile["supervisorId"] = user.SupervisorId
	profile["supervisorEmail"] = user.SupervisorEmail
	profile["user_id"] = user.Id

	return profile
}
