package alarms

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/spf13/cobra"
)

func InitCreateAlarmCommand(cwClient *cloudwatch.Client, cloudWatchCmd *cobra.Command) {
	var alarmName, metricName, namespace, comparisonOperator string
	var threshold float64
	var evaluationPeriods int32
	var tags []string

	var createAlarmCmd = &cobra.Command{
		Use:   "create-alarm",
		Short: "Creates a CloudWatch alarm",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createAlarm(cwClient, alarmName, metricName, namespace, comparisonOperator, threshold, evaluationPeriods, tags)
		},
	}

	createAlarmCmd.Flags().StringVarP(&alarmName, "alarm-name", "n", "", "Name of the alarm")
	createAlarmCmd.Flags().StringVarP(&metricName, "metric-name", "m", "", "Name of the metric")
	createAlarmCmd.Flags().StringVarP(&namespace, "namespace", "s", "", "Namespace of the metric")
	createAlarmCmd.Flags().StringVarP(&comparisonOperator, "comparison-operator", "c", "", "Comparison operator for the alarm")
	createAlarmCmd.Flags().Float64VarP(&threshold, "threshold", "t", 0, "Threshold for the alarm")
	createAlarmCmd.Flags().Int32VarP(&evaluationPeriods, "evaluation-periods", "e", 1, "Number of evaluation periods")
	createAlarmCmd.Flags().StringSliceVarP(&tags, "tags", "g", []string{}, "Tags for the alarm (key=value)")

	cloudWatchCmd.AddCommand(createAlarmCmd)
}

func createAlarm(cwClient *cloudwatch.Client, alarmName, metricName, namespace, comparisonOperator string, threshold float64, evaluationPeriods int32, tags []string) error {
	tagList := []types.Tag{}
	for _, tag := range tags {
		parts := strings.SplitN(tag, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid tag format: %s", tag)
		}
		tagList = append(tagList, types.Tag{
			Key:   &parts[0],
			Value: &parts[1],
		})
	}

	input := &cloudwatch.PutMetricAlarmInput{
		AlarmName:          &alarmName,
		MetricName:         &metricName,
		Namespace:          &namespace,
		ComparisonOperator: types.ComparisonOperator(comparisonOperator),
		Threshold:          &threshold,
		EvaluationPeriods:  &evaluationPeriods,
		Tags:               tagList,
	}

	_, err := cwClient.PutMetricAlarm(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("could not create alarm: %w", err)
	}

	fmt.Printf("Created alarm %s\n", alarmName)
	return nil
}
