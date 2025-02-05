package backup

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/cloudquery/cloudquery/plugins/source/aws/client"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func planSelections() *schema.Table {
	tableName := "aws_backup_plan_selections"
	return &schema.Table{
		Name:        tableName,
		Description: `https://docs.aws.amazon.com/aws-backup/latest/devguide/API_GetBackupSelection.html`,
		Resolver:    fetchBackupPlanSelections,
		Multiplex:   client.ServiceAccountRegionMultiplexer(tableName, "backup"),
		Transform:   transformers.TransformWithStruct(&backup.GetBackupSelectionOutput{}),
		Columns: []schema.Column{
			client.DefaultAccountIDColumn(false),
			client.DefaultRegionColumn(false),
			{
				Name:     "plan_arn",
				Type:     schema.TypeString,
				Resolver: schema.ParentColumnResolver("arn"),
			},
		},
	}
}

func fetchBackupPlanSelections(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	plan := parent.Item.(*backup.GetBackupPlanOutput)
	cl := meta.(*client.Client)
	svc := cl.Services().Backup
	params := backup.ListBackupSelectionsInput{
		BackupPlanId: plan.BackupPlanId,
		MaxResults:   aws.Int32(1000), // maximum value from https://docs.aws.amazon.com/aws-backup/latest/devguide/API_ListBackupSelections.html
	}
	for {
		result, err := svc.ListBackupSelections(ctx, &params)
		if err != nil {
			return err
		}
		for _, m := range result.BackupSelectionsList {
			s, err := svc.GetBackupSelection(
				ctx,
				&backup.GetBackupSelectionInput{BackupPlanId: plan.BackupPlanId, SelectionId: m.SelectionId},
				func(o *backup.Options) {
					o.Region = cl.Region
				},
			)
			if err != nil {
				return err
			}
			if s != nil {
				res <- *s
			}
		}
		if aws.ToString(result.NextToken) == "" {
			break
		}
		params.NextToken = result.NextToken
	}
	return nil
}
