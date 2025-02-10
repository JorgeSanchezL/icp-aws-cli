package metrics

import (
	"context"
	"fmt"
	"icp-aws-cli/pkg/utils"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/spf13/cobra"
)

func InitDeleteMetricCommand(cwClient *cloudwatch.Client, cloudWatchCmd *cobra.Command) {
	var metricName string
	var prefix string
	var pattern string
	var namespace string
	var dimensionName string
	var dimensionValue string
	var allMetrics bool

	var deleteMetricCmd = &cobra.Command{
		Use:   "delete-metrics",
		Short: "Deletes CloudWatch metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			if allMetrics && (metricName != "" || prefix != "" || pattern != "" || namespace != "" || dimensionName != "" || dimensionValue != "") {
				return fmt.Errorf("the --all flag cannot be combined with other filters")
			}

			if allMetrics {
				if !utils.ConfirmAction() {
					return fmt.Errorf("action cancelled by user")
				}
				return deleteAllMetrics(cwClient)
			}

			if metricName != "" && (prefix != "" || pattern != "" || namespace != "" || dimensionName != "" || dimensionValue != "") {
				return fmt.Errorf("metric name cannot be combined with other filters")
			}

			if metricName != "" {
				return deleteMetricByName(cwClient, metricName)
			}

			if prefix == "" && pattern == "" && namespace == "" && dimensionName == "" && dimensionValue == "" {
				return fmt.Errorf("at least one filter must be specified")
			}

			return deleteMetricsWithFilters(cwClient, prefix, pattern, namespace, dimensionName, dimensionValue)
		},
	}

	deleteMetricCmd.Flags().StringVarP(&metricName, "metric-name", "n", "", "Metric name to filter metrics")
	deleteMetricCmd.Flags().StringVarP(&prefix, "prefix", "x", "", "Prefix to filter metrics by name")
	deleteMetricCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter metrics by name")
	deleteMetricCmd.Flags().StringVarP(&namespace, "namespace", "s", "", "Namespace to filter metrics")
	deleteMetricCmd.Flags().StringVarP(&dimensionName, "dimension-name", "d", "", "Dimension name to filter metrics")
	deleteMetricCmd.Flags().StringVarP(&dimensionValue, "dimension-value", "v", "", "Dimension value to filter metrics")
	deleteMetricCmd.Flags().BoolVarP(&allMetrics, "all", "a", false, "Delete all metrics")

	cloudWatchCmd.AddCommand(deleteMetricCmd)
}

func deleteAllMetrics(cwClient *cloudwatch.Client) error {
	result, err := cwClient.ListMetrics(context.TODO(), &cloudwatch.ListMetricsInput{})
	if err != nil {
		return fmt.Errorf("could not list metrics: %w", err)
	}

	for _, metric := range result.Metrics {
		_, err := cwClient.DeleteAlarms(context.TODO(), &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{*metric.MetricName},
		})
		if err != nil {
			return fmt.Errorf("could not delete metric %s: %w", *metric.MetricName, err)
		}
		fmt.Printf("Deleted metric %s\n", *metric.MetricName)
	}

	return nil
}

func deleteMetricByName(cwClient *cloudwatch.Client, metricName string) error {
	_, err := cwClient.DeleteAlarms(context.TODO(), &cloudwatch.DeleteAlarmsInput{
		AlarmNames: []string{metricName},
	})
	if err != nil {
		return fmt.Errorf("could not delete metric %s: %w", metricName, err)
	}

	fmt.Printf("Deleted metric %s\n", metricName)
	return nil
}

func deleteMetricsWithFilters(cwClient *cloudwatch.Client, prefix, pattern, namespace, dimensionName, dimensionValue string) error {
	input := &cloudwatch.ListMetricsInput{}

	if namespace != "" {
		input.Namespace = &namespace
	}

	if dimensionName != "" && dimensionValue != "" {
		input.Dimensions = []types.DimensionFilter{
			{
				Name:  &dimensionName,
				Value: &dimensionValue,
			},
		}
	}

	if prefix != "" {
		input.MetricName = &prefix
	}

	result, err := cwClient.ListMetrics(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("could not list metrics: %w", err)
	}

	var metrics []types.Metric
	if pattern != "" {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid pattern: %w", err)
		}
		for _, metric := range result.Metrics {
			if re.MatchString(*metric.MetricName) {
				metrics = append(metrics, metric)
			}
		}
	} else {
		metrics = result.Metrics
	}

	for _, metric := range metrics {
		_, err := cwClient.DeleteAlarms(context.TODO(), &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []string{*metric.MetricName},
		})
		if err != nil {
			return fmt.Errorf("could not delete metric %s: %w", *metric.MetricName, err)
		}
		fmt.Printf("Deleted metric %s\n", *metric.MetricName)
	}

	return nil
}
