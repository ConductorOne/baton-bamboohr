package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	CompanyDomainField = field.StringField(
		"company-domain",
		field.WithDisplayName("Company Domain"),
		field.WithDescription("The company domain for your BambooHR account"),
		field.WithRequired(true),
	)
	ApiKeyField = field.StringField(
		"api-key",
		field.WithDisplayName("API Key"),
		field.WithDescription("The api key for your BambooHR account"),
		field.WithRequired(true),
		field.WithIsSecret(true),
	)
	BaseURLField = field.StringField(
		"base-url",
		field.WithDescription("Override the BambooHR API URL (for testing)"),
		field.WithHidden(true),
		field.WithExportTarget(field.ExportTargetCLIOnly),
	)
	CustomFieldsField = field.StringSliceField(
		"custom-fields",
		field.WithDisplayName("Custom Fields"),
		field.WithDescription("Additional BambooHR field names (standard or custom) to include in synced user profiles, e.g. customField4444,nickname."),
	)

	ConfigurationFields = []field.SchemaField{
		CompanyDomainField,
		ApiKeyField,
		BaseURLField,
		CustomFieldsField,
	}
	Configuration = field.NewConfiguration(ConfigurationFields)

	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run ./gen
var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("Bamboo HR V2"),
	field.WithHelpUrl("/docs/baton/bamboohr-v2"),
	field.WithIconUrl("/static/app-icons/bamboohr.svg"),
)
