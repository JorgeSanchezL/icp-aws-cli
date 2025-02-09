package loggroups

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

func InitCreateLogGroupCommand(cwLogsClient *cloudwatchlogs.Client, cloudWatchCmd *cobra.Command) {
	var logGroupName string

	var createLogGroupCmd = &cobra.Command{
		Use:   "create-log-group",
		Short: "Creates a CloudWatch log group",
		RunE: func(cmd *cobra.Command, args []string) error {
			if logGroupName == "" {
				return fmt.Errorf("log group name must be specified")
			}
			return createLogGroup(cwLogsClient, logGroupName)
		},
	}

	createLogGroupCmd.Flags().StringVarP(&logGroupName, "log-group-name", "n", "", "Name of the log group to create")
	cloudWatchCmd.AddCommand(createLogGroupCmd)
}

func createLogGroup(cwLogsClient *cloudwatchlogs.Client, logGroupName string) error {
	_, err := cwLogsClient.CreateLogGroup(context.TODO(), &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: &logGroupName,
	})
	if err != nil {
		return fmt.Errorf("could not create log group %s: %w", logGroupName, err)
	}

	fmt.Printf("Created log group %s\n", logGroupName)
	return nil
}
