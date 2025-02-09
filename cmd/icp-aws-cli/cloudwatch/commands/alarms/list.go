package alarms

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/spf13/cobra"
)

func InitListAlarmsCommand(cwClient *cloudwatch.Client, cloudWatchCmd *cobra.Command) {
	var alarmName string
	var pattern string
	var tagKey string
	var tagValue string
	var allAlarms bool

	var listAlarmsCmd = &cobra.Command{
		Use:   "list-alarms",
		Short: "Lists CloudWatch alarms",
		RunE: func(cmd *cobra.Command, args []string) error {
			if allAlarms && (alarmName != "" || pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("the --all flag cannot be combined with other filters")
			}

			if allAlarms {
				return listAllAlarms(cwClient)
			}

			if alarmName != "" && (pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("alarm name cannot be combined with other filters")
			}

			if alarmName != "" {
				return listAlarmsByName(cwClient, alarmName)
			}

			if pattern == "" && tagKey == "" && tagValue == "" {
				return fmt.Errorf("at least one filter must be specified")
			}

			return listAlarmsWithFilters(cwClient, pattern, tagKey, tagValue)
		},
	}

	listAlarmsCmd.Flags().StringVarP(&alarmName, "alarm-name", "a", "", "Alarm name to filter alarms")
	listAlarmsCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter alarms by name")
	listAlarmsCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter alarms")
	listAlarmsCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter alarms")
	listAlarmsCmd.Flags().BoolVarP(&allAlarms, "all", "l", false, "List all alarms")

	cloudWatchCmd.AddCommand(listAlarmsCmd)
}

func listAllAlarms(cwClient *cloudwatch.Client) error {
	result, err := cwClient.DescribeAlarms(context.TODO(), &cloudwatch.DescribeAlarmsInput{})
	if err != nil {
		return fmt.Errorf("could not list alarms: %w", err)
	}

	for _, alarm := range result.MetricAlarms {
		printAlarm(alarm)
	}

	return nil
}

func listAlarmsByName(cwClient *cloudwatch.Client, alarmName string) error {
	result, err := cwClient.DescribeAlarms(context.TODO(), &cloudwatch.DescribeAlarmsInput{
		AlarmNames: []string{alarmName},
	})
	if err != nil {
		return fmt.Errorf("could not list alarms: %w", err)
	}

	for _, alarm := range result.MetricAlarms {
		printAlarm(alarm)
	}

	return nil
}

func listAlarmsWithFilters(cwClient *cloudwatch.Client, pattern, tagKey, tagValue string) error {
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
		printAlarm(alarm)
	}

	return nil
}

func printAlarm(alarm types.MetricAlarm) {
	fmt.Printf("Alarm: %s\n", *alarm.AlarmName)
	fmt.Printf("Metric: %s\n", *alarm.MetricName)
	fmt.Printf("Namespace: %s\n", *alarm.Namespace)
	fmt.Printf("State: %s\n", alarm.StateValue)
	fmt.Printf("Threshold: %f\n", *alarm.Threshold)
	fmt.Printf("Comparison Operator: %s\n", alarm.ComparisonOperator)
	fmt.Printf("Evaluation Periods: %d\n", *alarm.EvaluationPeriods)
}
