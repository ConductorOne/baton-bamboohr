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

func WithRateLimitAnnotations(
	ratelimitDescriptionAnnotations ...*v2.RateLimitDescription,
) annotations.Annotations {
	outputAnnotations := annotations.Annotations{}
	for _, annotation := range ratelimitDescriptionAnnotations {
		outputAnnotations.Append(annotation)
	}

	return outputAnnotations
}

func (o *UserResourceType) List(
	ctx context.Context,
	_ *v2.ResourceId,
	_ *pagination.Token,
) ([]*v2.Resource, string, annotations.Annotations, error) {
	users, ratelimitData, err := o.bambooHRClient.ListUsers(ctx)
	outputAnnotations := WithRateLimitAnnotations(ratelimitData)
	if err != nil {
		return nil, "", outputAnnotations, err
	}

	rv := make([]*v2.Resource, 0)
	for _, user := range users {
		newResource, err := userResource(user)
		if err != nil {
			return nil, "", nil, err
		}

		rv = append(rv, newResource)
	}

	return rv, "", outputAnnotations, nil
}

func (o *UserResourceType) Entitlements(
	_ context.Context,
	_ *v2.Resource,
	_ *pagination.Token,
) (
	[]*v2.Entitlement,
	string,
	annotations.Annotations,
	error,
) {
	return nil, "", nil, nil
}

func (o *UserResourceType) Grants(
	_ context.Context,
	_ *v2.Resource,
	_ *pagination.Token,
) (
	[]*v2.Grant,
	string,
	annotations.Annotations,
	error,
) {
	return nil, "", nil, nil
}

func userBuilder(bambooHRClient *client.BambooHRClient) *UserResourceType {
	return &UserResourceType{
		resourceType:   resourceTypeUser,
		bambooHRClient: bambooHRClient,
	}
}

// userResource convert a BambooHR into a Resource.
func userResource(user *client.User) (*v2.Resource, error) {
	var userStatus v2.UserTrait_Status_Status
	profile := userProfile(user)

	displayName := fmt.Sprintf(
		"%s %s",
		user.FirstName,
		user.LastName,
	)

	switch user.Status {
	case "Active":
		userStatus = v2.UserTrait_Status_STATUS_ENABLED
	case "Inactive":
		userStatus = v2.UserTrait_Status_STATUS_DISABLED
	default:
		userStatus = v2.UserTrait_Status_STATUS_UNSPECIFIED
	}

	userTraitOptions := []resource.UserTraitOption{
		resource.WithUserProfile(profile),
		resource.WithEmail(user.Email, true),
		resource.WithStatus(userStatus),
	}

	return resource.NewUserResource(
		displayName,
		resourceTypeUser,
		user.Id,
		userTraitOptions,
	)
}

func userProfile(user *client.User) map[string]interface{} {
	profile := map[string]interface{}{
		"supervisorEId":      user.SupervisorEId,
		"supervisorFullName": user.Supervisor,
		"supervisorId":       user.SupervisorId,
		"supervisorEmail":    user.SupervisorEmail,
		"user_id":            user.Id,
		"firstName":          user.FirstName,
		"lastName":           user.LastName,
		"division":           user.Division,
		"department":         user.Department,
		"status":             user.Status,
		"employmentStatus":   user.EmploymentHistoryStatus,
		"terminationDate":    user.TerminationDate,
		"jobTitle":           user.JobTitle,
		"hireDate":           user.HireDate,
		"originalHireDate":   user.OriginalHireDate,
		"employeeNumber":     user.EmployeeNumber,
		"location":           user.Location,
		"preferredName":      user.PreferredName,
		"displayName":        user.DisplayName,
		"reportsTo":          user.ReportsTo,
	}

	// Merge configured custom fields. Don't let a custom field clobber a
	// standard profile attribute.
	for k, v := range user.CustomFields {
		if _, exists := profile[k]; exists {
			continue
		}
		profile[k] = v
	}

	return profile
}
