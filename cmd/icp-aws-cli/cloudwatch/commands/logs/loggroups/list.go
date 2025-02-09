package loggroups

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/spf13/cobra"
)

func InitListLogGroupsCommand(cwClient *cloudwatchlogs.Client, cloudWatchCmd *cobra.Command) {
	var logGroupName string
	var pattern string
	var tagKey string
	var tagValue string
	var allLogs bool

	var listLogsCmd = &cobra.Command{
		Use:   "list-log-groups",
		Short: "Lists CloudWatch log groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			if allLogs && (logGroupName != "" || pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("the --all flag cannot be combined with other filters")
			}

			if allLogs {
				return listAllLogs(cwClient)
			}

			if logGroupName != "" && (pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("log group name cannot be combined with other filters")
			}

			if logGroupName != "" {
				return listLogsByName(cwClient, logGroupName)
			}

			if pattern == "" && tagKey == "" && tagValue == "" {
				return fmt.Errorf("at least one filter must be specified")
			}

			return listLogsWithFilters(cwClient, pattern, tagKey, tagValue)
		},
	}

	listLogsCmd.Flags().StringVarP(&logGroupName, "log-group-name", "n", "", "Log group name to filter logs")
	listLogsCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter logs by name")
	listLogsCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter logs")
	listLogsCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter logs")
	listLogsCmd.Flags().BoolVarP(&allLogs, "all", "a", false, "List all logs")

	cloudWatchCmd.AddCommand(listLogsCmd)
}

func listAllLogs(cwClient *cloudwatchlogs.Client) error {
	result, err := cwClient.DescribeLogGroups(context.TODO(), &cloudwatchlogs.DescribeLogGroupsInput{})
	if err != nil {
		return fmt.Errorf("could not list log groups: %w", err)
	}

	for _, logGroup := range result.LogGroups {
		printLogGroup(logGroup)
	}

	return nil
}

func listLogsByName(cwClient *cloudwatchlogs.Client, logGroupName string) error {
	result, err := cwClient.DescribeLogGroups(context.TODO(), &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: &logGroupName,
	})
	if err != nil {
		return fmt.Errorf("could not list log groups: %w", err)
	}

	for _, logGroup := range result.LogGroups {
		printLogGroup(logGroup)
	}

	return nil
}

func listLogsWithFilters(cwClient *cloudwatchlogs.Client, pattern, tagKey, tagValue string) error {
	result, err := cwClient.DescribeLogGroups(context.TODO(), &cloudwatchlogs.DescribeLogGroupsInput{})
	if err != nil {
		return fmt.Errorf("could not list log groups: %w", err)
	}

	var logGroups []types.LogGroup
	if pattern != "" {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid pattern: %w", err)
		}
		for _, logGroup := range result.LogGroups {
			if re.MatchString(*logGroup.LogGroupName) {
				logGroups = append(logGroups, logGroup)
			}
		}
	} else {
		logGroups = result.LogGroups
	}

	if tagKey != "" && tagValue != "" {
		var filteredLogGroups []types.LogGroup
		for _, logGroup := range logGroups {
			tags, err := cwClient.ListTagsForResource(context.TODO(), &cloudwatchlogs.ListTagsForResourceInput{
				ResourceArn: logGroup.Arn,
			})
			if err != nil {
				return fmt.Errorf("could not list tags for log group %s: %w", *logGroup.LogGroupName, err)
			}
			for key, value := range tags.Tags {
				if key == tagKey && value == tagValue {
					filteredLogGroups = append(filteredLogGroups, logGroup)
					break
				}
			}
		}
		logGroups = filteredLogGroups
	}

	for _, logGroup := range logGroups {
		printLogGroup(logGroup)
	}

	return nil
}

func printLogGroup(logGroup types.LogGroup) {
	creationTime := time.Unix(0, *logGroup.CreationTime*int64(time.Millisecond)).Format("2006-01-02 15:04:05")
	fmt.Printf("Log Group Name: %s, Creation Time: %s\n", *logGroup.LogGroupName, creationTime)
}
