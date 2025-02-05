package lightsail

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"
	"github.com/cloudquery/cloudquery/plugins/source/aws/client"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func databaseEvents() *schema.Table {
	tableName := "aws_lightsail_database_events"
	return &schema.Table{
		Name:        tableName,
		Description: `https://docs.aws.amazon.com/lightsail/2016-11-28/api-reference/API_RelationalDatabaseEvent.html`,
		Resolver:    fetchLightsailDatabaseEvents,
		Transform:   transformers.TransformWithStruct(&types.RelationalDatabaseEvent{}),
		Multiplex:   client.ServiceAccountRegionMultiplexer(tableName, "lightsail"),
		Columns: []schema.Column{
			client.DefaultAccountIDColumn(false),
			client.DefaultRegionColumn(false),
			{
				Name:     "database_arn",
				Type:     schema.TypeString,
				Resolver: schema.ParentColumnResolver("arn"),
			},
		},
	}
}

func fetchLightsailDatabaseEvents(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	r := parent.Item.(types.RelationalDatabase)
	input := lightsail.GetRelationalDatabaseEventsInput{
		RelationalDatabaseName: r.Name,
		DurationInMinutes:      aws.Int32(20160), // two weeks
	}
	c := meta.(*client.Client)
	svc := c.Services().Lightsail
	for {
		response, err := svc.GetRelationalDatabaseEvents(ctx, &input)
		if err != nil {
			return err
		}
		res <- response.RelationalDatabaseEvents
		if aws.ToString(response.NextPageToken) == "" {
			break
		}
		input.PageToken = response.NextPageToken
	}
	return nil
}
