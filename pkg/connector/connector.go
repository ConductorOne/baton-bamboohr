package connector

import (
	"context"
	"fmt"
	"io"

	"github.com/conductorone/baton-bamboohr/pkg/connector/client"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

type BambooHr struct {
	customerDomain string
	client         *client.BambooHRClient
	apiKey         string
}

func New(
	ctx context.Context,
	customerDomain string,
	apiKey string,
) (*BambooHr, error) {
	client, err := client.New(ctx, apiKey, customerDomain)
	if err != nil {
		return nil, err
	}
	rv := &BambooHr{
		customerDomain: customerDomain,
		apiKey:         apiKey,
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
	return []connectorbuilder.ResourceSyncer{
		userBuilder(c.client),
	}
}
