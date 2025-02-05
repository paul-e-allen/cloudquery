package frauddetector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/frauddetector"
	"github.com/aws/aws-sdk-go-v2/service/frauddetector/types"
	"github.com/cloudquery/cloudquery/plugins/source/aws/client"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func modelVersions() *schema.Table {
	return &schema.Table{
		Name:        "aws_frauddetector_model_versions",
		Description: `https://docs.aws.amazon.com/frauddetector/latest/api/API_ModelVersionDetail.html`,
		Resolver:    fetchFrauddetectorModelVersions,
		Transform:   transformers.TransformWithStruct(&types.ModelVersionDetail{}),
		Columns: []schema.Column{
			client.DefaultAccountIDColumn(false),
			client.DefaultRegionColumn(false),
			{
				Name:     "arn",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("Arn"),
				CreationOptions: schema.ColumnCreationOptions{
					PrimaryKey: true,
				},
			},
		},
	}
}

func fetchFrauddetectorModelVersions(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	paginator := frauddetector.NewDescribeModelVersionsPaginator(meta.(*client.Client).Services().Frauddetector,
		&frauddetector.DescribeModelVersionsInput{ModelId: parent.Item.(types.Model).ModelId})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		res <- output.ModelVersionDetails
	}
	return nil
}
