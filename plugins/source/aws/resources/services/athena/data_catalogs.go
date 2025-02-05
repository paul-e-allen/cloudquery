package athena

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
	"github.com/cloudquery/cloudquery/plugins/source/aws/client"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func DataCatalogs() *schema.Table {
	tableName := "aws_athena_data_catalogs"
	return &schema.Table{
		Name:                tableName,
		Description:         `https://docs.aws.amazon.com/athena/latest/APIReference/API_DataCatalog.html`,
		Resolver:            fetchAthenaDataCatalogs,
		PreResourceResolver: getDataCatalog,
		Multiplex:           client.ServiceAccountRegionMultiplexer(tableName, "athena"),
		Transform:           transformers.TransformWithStruct(&types.DataCatalog{}),
		Columns: []schema.Column{
			client.DefaultAccountIDColumn(false),
			client.DefaultRegionColumn(false),
			{
				Name:     "arn",
				Type:     schema.TypeString,
				Resolver: resolveAthenaDataCatalogArn,
				CreationOptions: schema.ColumnCreationOptions{
					PrimaryKey: true,
				},
			},
			{
				Name:     "tags",
				Type:     schema.TypeJSON,
				Resolver: resolveAthenaDataCatalogTags,
			},
		},

		Relations: []*schema.Table{
			dataCatalogDatabases(),
		},
	}
}

func fetchAthenaDataCatalogs(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	c := meta.(*client.Client)
	svc := c.Services().Athena
	input := athena.ListDataCatalogsInput{}
	for {
		response, err := svc.ListDataCatalogs(ctx, &input)
		if err != nil {
			return err
		}
		res <- response.DataCatalogsSummary
		if aws.ToString(response.NextToken) == "" {
			break
		}
		input.NextToken = response.NextToken
	}
	return nil
}

func getDataCatalog(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource) error {
	c := meta.(*client.Client)
	svc := c.Services().Athena
	catalogSummary := resource.Item.(types.DataCatalogSummary)
	dc, err := svc.GetDataCatalog(ctx, &athena.GetDataCatalogInput{
		Name: catalogSummary.CatalogName,
	})
	if err != nil {
		// retrieving of default data catalog (AwsDataCatalog) returns "not found error" (with statuscode 400: InvalidRequestException:...) but it exists and its
		// relations can be fetched by its name
		if client.IsAWSError(err, "InvalidRequestException") && *catalogSummary.CatalogName == "AwsDataCatalog" {
			resource.Item = types.DataCatalog{Name: catalogSummary.CatalogName, Type: catalogSummary.Type}
			return nil
		}
		return err
	}
	resource.Item = *dc.DataCatalog
	return nil
}

func resolveAthenaDataCatalogArn(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	cl := meta.(*client.Client)
	dc := resource.Item.(types.DataCatalog)
	return resource.Set(c.Name, createDataCatalogArn(cl, *dc.Name))
}

func resolveAthenaDataCatalogTags(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	cl := meta.(*client.Client)
	svc := cl.Services().Athena
	dc := resource.Item.(types.DataCatalog)
	arnStr := createDataCatalogArn(cl, *dc.Name)
	paginator := athena.NewListTagsForResourcePaginator(svc, &athena.ListTagsForResourceInput{ResourceARN: &arnStr})
	tags := make(map[string]string)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if cl.IsNotFoundError(err) {
				return nil
			}
			return err
		}
		client.TagsIntoMap(page.Tags, tags)
	}
	return resource.Set(c.Name, tags)
}

func createDataCatalogArn(cl *client.Client, catalogName string) string {
	return arn.ARN{
		Partition: cl.Partition,
		Service:   string(client.Athena),
		Region:    cl.Region,
		AccountID: cl.AccountID,
		Resource:  fmt.Sprintf("datacatalog/%s", catalogName),
	}.String()
}
