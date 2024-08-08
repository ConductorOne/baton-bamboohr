package main

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
	configurationFields = []field.SchemaField{
		CompanyDomainField,
		ApiKeyField,
	}
	Configuration = field.NewConfiguration(configurationFields)
)
