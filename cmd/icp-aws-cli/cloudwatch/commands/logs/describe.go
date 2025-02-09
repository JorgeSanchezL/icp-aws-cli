package logs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

func InitDescribeLogGroupCommand(cwClient *cloudwatchlogs.Client, cloudWatchCmd *cobra.Command) {
	var logGroupName string

	var describeLogCmd = &cobra.Command{
		Use:   "describe-log-group",
		Short: "Describes a CloudWatch log group",
		RunE: func(cmd *cobra.Command, args []string) error {
			return describeLog(cwClient, logGroupName)
		},
	}

	describeLogCmd.Flags().StringVarP(&logGroupName, "log-group-name", "l", "", "Name of the log group")
	cloudWatchCmd.AddCommand(describeLogCmd)
}

func describeLog(cwClient *cloudwatchlogs.Client, logGroupName string) error {
	if logGroupName == "" {
		return fmt.Errorf("log group name must be specified")
	}

	result, err := cwClient.DescribeLogGroups(context.TODO(), &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: &logGroupName,
	})
	if err != nil {
		return fmt.Errorf("could not describe log group: %w", err)
	}

	for _, logGroup := range result.LogGroups {
		fmt.Printf("Log Group: %s\n", *logGroup.LogGroupName)
		fmt.Printf("Creation Time: %d\n", logGroup.CreationTime)
		fmt.Printf("Stored Bytes: %d\n", logGroup.StoredBytes)
		fmt.Printf("Retention In Days: %d\n", logGroup.RetentionInDays)
	}

	return nil
}
