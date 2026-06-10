package main

import (
	"context"
	"fmt"
	"os"

	cfg "github.com/conductorone/baton-bamboohr/pkg/config"
	"github.com/conductorone/baton-bamboohr/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var version = "dev"

func main() {
	ctx := context.Background()

	_, cmd, err := config.DefineConfiguration(
		ctx,
		"baton-bamboohr",
		getConnector,
		cfg.Configuration,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, bhrc *cfg.Bamboohr) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	err := field.Validate(cfg.Config, bhrc)
	if err != nil {
		return nil, err
	}

	companyDomain := bhrc.CompanyDomain
	if companyDomain == "" {
		return nil, fmt.Errorf("company domain field is required")
	}

	apiKey := bhrc.ApiKey
	if apiKey == "" {
		return nil, fmt.Errorf("api key field is required")
	}

	cb, err := connector.New(
		ctx,
		companyDomain,
		apiKey,
		bhrc.BaseUrl,
		bhrc.CustomFields,
	)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	conn, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}
	return conn, nil
}
