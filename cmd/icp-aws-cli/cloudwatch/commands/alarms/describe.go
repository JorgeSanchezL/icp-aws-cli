package alarms

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/spf13/cobra"
)

func InitDescribeAlarmCommand(cwClient *cloudwatch.Client, cloudWatchCmd *cobra.Command) {
	var alarmName string

	var describeAlarmCmd = &cobra.Command{
		Use:   "describe-alarm",
		Short: "Describes a CloudWatch alarm",
		RunE: func(cmd *cobra.Command, args []string) error {
			return describeAlarm(cwClient, alarmName)
		},
	}

	describeAlarmCmd.Flags().StringVarP(&alarmName, "alarm-name", "a", "", "Name of the alarm")
	cloudWatchCmd.AddCommand(describeAlarmCmd)
}

func describeAlarm(cwClient *cloudwatch.Client, alarmName string) error {
	result, err := cwClient.DescribeAlarms(context.TODO(), &cloudwatch.DescribeAlarmsInput{
		AlarmNames: []string{alarmName},
	})
	if err != nil {
		return fmt.Errorf("could not describe alarm: %w", err)
	}

	for _, alarm := range result.MetricAlarms {
		fmt.Printf("Alarm: %s\n", *alarm.AlarmName)
		fmt.Printf("Metric: %s\n", *alarm.MetricName)
		fmt.Printf("Namespace: %s\n", *alarm.Namespace)
		fmt.Printf("State: %s\n", alarm.StateValue)
		fmt.Printf("Threshold: %f\n", *alarm.Threshold)
		fmt.Printf("Comparison Operator: %s\n", alarm.ComparisonOperator)
		fmt.Printf("Evaluation Periods: %d\n", *alarm.EvaluationPeriods)
	}

	return nil
}
