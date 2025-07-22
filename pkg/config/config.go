package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	CompanyDomainField = field.StringField(
		"company-domain",
		field.WithDescription("The company domain for your BambooHR account"),
		field.WithRequired(true),
	)
	ApiKeyField = field.StringField(
		"api-key",
		field.WithDescription("The api key for your BambooHR account"),
		field.WithRequired(true),
	)

	ConfigurationFields = []field.SchemaField{
		CompanyDomainField,
		ApiKeyField,
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
