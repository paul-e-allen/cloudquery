package client

import "github.com/cloudquery/plugin-sdk/schema"

var OrganizationID = schema.Column{
	Name:     "organization_id",
	Type:     schema.TypeString,
	Resolver: ResolveOrganizationID,
	CreationOptions: schema.ColumnCreationOptions{
		PrimaryKey: true,
		NotNull:    true,
	},
}
