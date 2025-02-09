package streams

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/spf13/cobra"
)

func InitListLogStreamsCommand(cwLogsClient *cloudwatchlogs.Client, cloudWatchCmd *cobra.Command) {
	var logGroupName string
	var limit int32

	var listLogStreamsCmd = &cobra.Command{
		Use:   "list-log-streams",
		Short: "Lists the last CloudWatch log streams in a log group (10 by default)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if logGroupName == "" {
				return fmt.Errorf("log group name must be specified")
			}
			return listLogStreams(cwLogsClient, logGroupName, limit)
		},
	}

	listLogStreamsCmd.Flags().StringVarP(&logGroupName, "log-group-name", "n", "", "Name of the log group")
	listLogStreamsCmd.Flags().Int32VarP(&limit, "limit", "l", 10, "Number of log streams to list")
	cloudWatchCmd.AddCommand(listLogStreamsCmd)
}

func listLogStreams(cwLogsClient *cloudwatchlogs.Client, logGroupName string, limit int32) error {
	result, err := cwLogsClient.DescribeLogStreams(context.TODO(), &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: &logGroupName,
		Limit:        &limit,
	})
	if err != nil {
		return fmt.Errorf("could not list log streams: %w", err)
	}

	for _, logStream := range result.LogStreams {
		printLogStream(logStream)
	}

	return nil
}

func printLogStream(logStream types.LogStream) {
	creationTime := time.Unix(0, *logStream.CreationTime*int64(time.Millisecond)).Format("2006-01-02 15:04:05")
	fmt.Printf("Log Stream Name: %s, Creation Time: %s, Last Event Time: %s\n", *logStream.LogStreamName, creationTime, time.Unix(0, *logStream.LastEventTimestamp*int64(time.Millisecond)).Format("2006-01-02 15:04:05"))
}
