package events

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/spf13/cobra"
)

func InitGetLogEventsCommand(cwLogsClient *cloudwatchlogs.Client, cloudWatchCmd *cobra.Command) {
	var logGroupName string
	var logStreamName string
	var limit int32

	var getLogEventsCmd = &cobra.Command{
		Use:   "get-log-events",
		Short: "Gets the last log events from a CloudWatch log stream (10 by default)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if logGroupName == "" || logStreamName == "" {
				return fmt.Errorf("log group name and log stream name must be specified")
			}
			return getLogEvents(cwLogsClient, logGroupName, logStreamName, limit)
		},
	}

	getLogEventsCmd.Flags().StringVarP(&logGroupName, "log-group-name", "g", "", "Name of the log group")
	getLogEventsCmd.Flags().StringVarP(&logStreamName, "log-stream-name", "s", "", "Name of the log stream")
	getLogEventsCmd.Flags().Int32VarP(&limit, "limit", "l", 10, "Number of log events to get")
	cloudWatchCmd.AddCommand(getLogEventsCmd)
}

func getLogEvents(cwLogsClient *cloudwatchlogs.Client, logGroupName, logStreamName string, limit int32) error {
	result, err := cwLogsClient.GetLogEvents(context.TODO(), &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  &logGroupName,
		LogStreamName: &logStreamName,
		Limit:         &limit,
	})
	if err != nil {
		return fmt.Errorf("could not get log events: %w", err)
	}

	for _, event := range result.Events {
		printLogEvent(event)
	}

	return nil
}

func printLogEvent(event types.OutputLogEvent) {
	timestamp := time.Unix(0, *event.Timestamp*int64(time.Millisecond)).Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] %s\n", timestamp, *event.Message)
}
