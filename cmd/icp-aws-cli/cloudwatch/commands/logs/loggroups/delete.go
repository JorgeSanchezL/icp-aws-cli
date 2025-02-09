package loggroups

import (
	"context"
	"fmt"
	"icp-aws-cli/pkg/utils"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/spf13/cobra"
)

func InitDeleteLogGroupCommand(cwLogsClient *cloudwatchlogs.Client, cloudWatchCmd *cobra.Command) {
	var logGroupName string
	var prefix string
	var pattern string
	var tagKey string
	var tagValue string
	var allLogGroups bool

	var deleteLogGroupCmd = &cobra.Command{
		Use:   "delete-loggroup",
		Short: "Deletes a CloudWatch log group",
		RunE: func(cmd *cobra.Command, args []string) error {
			if allLogGroups && (logGroupName != "" || prefix != "" || pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("the --all flag cannot be combined with other filters")
			}

			if allLogGroups {
				if !utils.ConfirmAction() {
					return fmt.Errorf("action cancelled by user")
				}
				return deleteAllLogGroups(cwLogsClient)
			}

			if logGroupName != "" && (prefix != "" || pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("log group name cannot be combined with other filters")
			}

			if logGroupName != "" {
				return deleteLogGroupByName(cwLogsClient, logGroupName)
			}

			if prefix == "" && pattern == "" && tagKey == "" && tagValue == "" {
				return fmt.Errorf("at least one filter must be specified")
			}

			return deleteLogGroupsWithFilters(cwLogsClient, prefix, pattern, tagKey, tagValue)
		},
	}

	deleteLogGroupCmd.Flags().StringVarP(&logGroupName, "log-group-name", "n", "", "Log group name to filter log groups")
	deleteLogGroupCmd.Flags().StringVarP(&prefix, "prefix", "x", "", "Prefix to filter log groups by name")
	deleteLogGroupCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter log groups by name")
	deleteLogGroupCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter log groups")
	deleteLogGroupCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter log groups")
	deleteLogGroupCmd.Flags().BoolVarP(&allLogGroups, "all", "a", false, "Delete all log groups")

	cloudWatchCmd.AddCommand(deleteLogGroupCmd)
}

func deleteAllLogGroups(cwLogsClient *cloudwatchlogs.Client) error {
	result, err := cwLogsClient.DescribeLogGroups(context.TODO(), &cloudwatchlogs.DescribeLogGroupsInput{})
	if err != nil {
		return fmt.Errorf("could not list log groups: %w", err)
	}

	for _, logGroup := range result.LogGroups {
		_, err := cwLogsClient.DeleteLogGroup(context.TODO(), &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: logGroup.LogGroupName,
		})
		if err != nil {
			return fmt.Errorf("could not delete log group %s: %w", *logGroup.LogGroupName, err)
		}
		fmt.Printf("Deleted log group %s\n", *logGroup.LogGroupName)
	}

	return nil
}

func deleteLogGroupByName(cwLogsClient *cloudwatchlogs.Client, logGroupName string) error {
	_, err := cwLogsClient.DeleteLogGroup(context.TODO(), &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: &logGroupName,
	})
	if err != nil {
		return fmt.Errorf("could not delete log group %s: %w", logGroupName, err)
	}

	fmt.Printf("Deleted log group %s\n", logGroupName)
	return nil
}

func deleteLogGroupsWithFilters(cwLogsClient *cloudwatchlogs.Client, prefix, pattern, tagKey, tagValue string) error {
	input := &cloudwatchlogs.DescribeLogGroupsInput{}

	if prefix != "" {
		input.LogGroupNamePrefix = &prefix
	}

	result, err := cwLogsClient.DescribeLogGroups(context.TODO(), input)
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
			tags, err := cwLogsClient.ListTagsLogGroup(context.TODO(), &cloudwatchlogs.ListTagsLogGroupInput{
				LogGroupName: logGroup.LogGroupName,
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
		_, err := cwLogsClient.DeleteLogGroup(context.TODO(), &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: logGroup.LogGroupName,
		})
		if err != nil {
			return fmt.Errorf("could not delete log group %s: %w", *logGroup.LogGroupName, err)
		}
		fmt.Printf("Deleted log group %s\n", *logGroup.LogGroupName)
	}

	return nil
}
