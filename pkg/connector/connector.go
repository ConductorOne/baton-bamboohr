package connector

import (
	"context"
	"fmt"
	"io"

	"github.com/ConductorOne/BambooHR/pkg/connector/client"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

var (
	resourceTypeUser = &v2.ResourceType{
		Id:          "user",
		DisplayName: "User",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_USER,
		},
		Annotations: v1AnnotationsForResourceType("user"),
	}
)

type Config struct {
	CompanyDomain string
	ApiKey        string
}

type BambooHr struct {
	customerDomain string
	client         *client.BambooHRClient
	apiKey         string
}

func New(ctx context.Context, config Config) (*BambooHr, error) {
	client, err := client.New(ctx, config.ApiKey, config.CompanyDomain)
	if err != nil {
		return nil, err
	}
	rv := &BambooHr{
		customerDomain: config.CompanyDomain,
		apiKey:         config.ApiKey,
		client:         client,
	}
	return rv, nil
}

func (c *BambooHr) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	_, err := c.Validate(ctx)
	if err != nil {
		return nil, err
	}

	var annos annotations.Annotations
	annos.Update(&v2.ExternalLink{
		Url: c.customerDomain,
	})

	return &v2.ConnectorMetadata{
		DisplayName: "BambooHR",
		Annotations: annos,
	}, nil
}

func (c *BambooHr) Validate(ctx context.Context) (annotations.Annotations, error) {
	_, err := c.client.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to validate API keys: %w", err)
	}
	return nil, nil
}

func (c *BambooHr) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

func (c *BambooHr) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	rs := []connectorbuilder.ResourceSyncer{}
	rs = append(rs, userBuilder(c.client))
	return rs
}
