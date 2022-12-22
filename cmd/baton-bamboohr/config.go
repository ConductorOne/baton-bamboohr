package main

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/spf13/cobra"
)

// config defines the external configuration required for the connector to run.
type config struct {
	cli.BaseConfig `mapstructure:",squash"` // Puts the base config options in the same place as the connector options

	CompanyDomain string `mapstructure:"company-domain"`
	ApiKey        string `mapstructure:"api-key"`
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	if cfg.CompanyDomain == "" {
		return fmt.Errorf("company domain is missing")
	}
	if cfg.ApiKey == "" {
		return fmt.Errorf("api key is missing")
	}
	return nil
}

// cmdFlags sets the cmdFlags required for the connector.
func cmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("company-domain", "", "The company domain for your BambooHR account. ($BATON_COMPANY_DOMAIN)")
	cmd.PersistentFlags().String("api-key", "", "The api key for your BambooHR account. ($BATON_API_KEY)")
}
