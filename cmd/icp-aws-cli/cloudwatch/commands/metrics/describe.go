package metrics

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/spf13/cobra"
)

func InitDescribeMetricCommand(cwClient *cloudwatch.Client, cloudWatchCmd *cobra.Command) {
	var metricName string
	var namespace string
	var dimensionName string
	var dimensionValue string

	var describeMetricCmd = &cobra.Command{
		Use:   "describe-metric",
		Short: "Describes a CloudWatch metric",
		RunE: func(cmd *cobra.Command, args []string) error {
			return describeMetric(cwClient, metricName, namespace, dimensionName, dimensionValue)
		},
	}

	describeMetricCmd.Flags().StringVarP(&metricName, "metric-name", "m", "", "Name of the metric")
	describeMetricCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace of the metric")
	describeMetricCmd.Flags().StringVarP(&dimensionName, "dimension-name", "d", "", "Dimension name of the metric")
	describeMetricCmd.Flags().StringVarP(&dimensionValue, "dimension-value", "v", "", "Dimension value of the metric")
	cloudWatchCmd.AddCommand(describeMetricCmd)
}

func describeMetric(cwClient *cloudwatch.Client, metricName, namespace, dimensionName, dimensionValue string) error {
	if metricName == "" || namespace == "" {
		return fmt.Errorf("metric name and namespace must be specified")
	}

	input := &cloudwatch.ListMetricsInput{
		MetricName: &metricName,
		Namespace:  &namespace,
	}

	if dimensionName != "" && dimensionValue != "" {
		input.Dimensions = []types.DimensionFilter{
			{
				Name:  &dimensionName,
				Value: &dimensionValue,
			},
		}
	}

	result, err := cwClient.ListMetrics(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("could not describe metric: %w", err)
	}

	for _, metric := range result.Metrics {
		fmt.Printf("Metric: %s\n", *metric.MetricName)
		for _, dimension := range metric.Dimensions {
			fmt.Printf("  Dimension: %s = %s\n", *dimension.Name, *dimension.Value)
		}
	}

	return nil
}
