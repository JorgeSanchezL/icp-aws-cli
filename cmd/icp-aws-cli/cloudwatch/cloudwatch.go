package cloudwatch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/spf13/cobra"
)

func InitCommands(cloudwatchClient *cloudwatch.Client) *cobra.Command {
	var cloudwatchCmd = &cobra.Command{
		Use:   "cloudwatch",
		Short: "Commands to interact with Amazon CloudWatch",
		Long:  "Allows listing and managing CloudWatch metrics in Amazon CloudWatch.",
	}

	var listMetricsCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists CloudWatch metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listMetrics(cloudwatchClient)
		},
	}

	cloudwatchCmd.AddCommand(listMetricsCmd)
	return cloudwatchCmd
}

func listMetrics(cloudwatchClient *cloudwatch.Client) error {
	result, err := cloudwatchClient.ListMetrics(context.TODO(), &cloudwatch.ListMetricsInput{})
	if err != nil {
		return fmt.Errorf("error listing CloudWatch metrics: %w", err)
	}

	for _, metric := range result.Metrics {
		fmt.Println(*metric.MetricName)
	}
	return nil
}
