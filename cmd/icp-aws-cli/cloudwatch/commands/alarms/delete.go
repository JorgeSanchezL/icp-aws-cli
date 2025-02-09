package alarms

import (
	"context"
	"fmt"
	"icp-aws-cli/pkg/utils"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/spf13/cobra"
)

func InitDeleteAlarmCommand(cwClient *cloudwatch.Client, cloudWatchCmd *cobra.Command) {
	var alarmName string
	var pattern string
	var tagKey string
	var tagValue string
	var allAlarms bool

	var deleteAlarmCmd = &cobra.Command{
		Use:   "delete-alarm",
		Short: "Deletes a CloudWatch alarm",
		RunE: func(cmd *cobra.Command, args []string) error {
			if allAlarms && (alarmName != "" || pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("the --all flag cannot be combined with other filters")
			}

			if allAlarms {
				if !utils.ConfirmAction() {
					return fmt.Errorf("action cancelled by user")
				}
				return deleteAllAlarms(cwClient)
			}

			if alarmName != "" && (pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("alarm name cannot be combined with other filters")
			}

			if alarmName != "" {
				return deleteAlarmByName(cwClient, alarmName)
			}

			if pattern == "" && tagKey == "" && tagValue == "" {
				return fmt.Errorf("at least one filter must be specified")
			}

			return deleteAlarmsWithFilters(cwClient, pattern, tagKey, tagValue)
		},
	}

	deleteAlarmCmd.Flags().StringVarP(&alarmName, "alarm-name", "a", "", "Alarm name to filter alarms")
	deleteAlarmCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter alarms by name")
	deleteAlarmCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter alarms")
	deleteAlarmCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter alarms")
	deleteAlarmCmd.Flags().BoolVarP(&allAlarms, "all", "l", false, "Delete all alarms")

	cloudWatchCmd.AddCommand(deleteAlarmCmd)
}

func deleteAllAlarms(cwClient *cloudwatch.Client) error {
	result, err := cwClient.DescribeAlarms(context.TODO(), &cloudwatch.DescribeAlarmsInput{})
	if err != nil {
		return fmt.Errorf("could not list alarms: %w", err)
	}

	for _, alarm := range result.MetricAlarms {
		_, err := cwClient.DeleteAlarms(context.TODO(), &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{*alarm.AlarmName},
		})
		if err != nil {
			return fmt.Errorf("could not delete alarm %s: %w", *alarm.AlarmName, err)
		}
		fmt.Printf("Deleted alarm %s\n", *alarm.AlarmName)
	}

	return nil
}

func deleteAlarmByName(cwClient *cloudwatch.Client, alarmName string) error {
	_, err := cwClient.DeleteAlarms(context.TODO(), &cloudwatch.DeleteAlarmsInput{
		AlarmNames: []string{alarmName},
	})
	if err != nil {
		return fmt.Errorf("could not delete alarm %s: %w", alarmName, err)
	}

	fmt.Printf("Deleted alarm %s\n", alarmName)
	return nil
}

func deleteAlarmsWithFilters(cwClient *cloudwatch.Client, pattern, tagKey, tagValue string) error {
	result, err := cwClient.DescribeAlarms(context.TODO(), &cloudwatch.DescribeAlarmsInput{})
	if err != nil {
		return fmt.Errorf("could not list alarms: %w", err)
	}

	var alarms []types.MetricAlarm
	if pattern != "" {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid pattern: %w", err)
		}
		for _, alarm := range result.MetricAlarms {
			if re.MatchString(*alarm.AlarmName) {
				alarms = append(alarms, alarm)
			}
		}
	} else {
		alarms = result.MetricAlarms
	}

	if tagKey != "" && tagValue != "" {
		var filteredAlarms []types.MetricAlarm
		for _, alarm := range alarms {
			tags, err := cwClient.ListTagsForResource(context.TODO(), &cloudwatch.ListTagsForResourceInput{
				ResourceARN: alarm.AlarmArn,
			})
			if err != nil {
				return fmt.Errorf("could not list tags for alarm %s: %w", *alarm.AlarmName, err)
			}
			for _, tag := range tags.Tags {
				if *tag.Key == tagKey && *tag.Value == tagValue {
					filteredAlarms = append(filteredAlarms, alarm)
					break
				}
			}
		}
		alarms = filteredAlarms
	}

	for _, alarm := range alarms {
		_, err := cwClient.DeleteAlarms(context.TODO(), &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{*alarm.AlarmName},
		})
		if err != nil {
			return fmt.Errorf("could not delete alarm %s: %w", *alarm.AlarmName, err)
		}
		fmt.Printf("Deleted alarm %s\n", *alarm.AlarmName)
	}

	return nil
}
