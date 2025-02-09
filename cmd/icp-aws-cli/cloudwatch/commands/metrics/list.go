package metrics

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/spf13/cobra"
)

func InitListMetricsCommand(cwClient *cloudwatch.Client, cloudWatchCmd *cobra.Command) {
	var metricName string
	var prefix string
	var pattern string
	var namespace string
	var dimensionName string
	var dimensionValue string
	var allMetrics bool

	var listMetricsCmd = &cobra.Command{
		Use:   "list-metrics",
		Short: "Lists CloudWatch metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			if allMetrics && (metricName != "" || prefix != "" || pattern != "" || namespace != "" || dimensionName != "" || dimensionValue != "") {
				return fmt.Errorf("the --all flag cannot be combined with other filters")
			}

			if allMetrics {
				return listAllMetrics(cwClient)
			}

			if metricName != "" && (prefix != "" || pattern != "" || namespace != "" || dimensionName != "" || dimensionValue != "") {
				return fmt.Errorf("metric name cannot be combined with other filters")
			}

			if metricName != "" {
				return listMetricsByName(cwClient, metricName)
			}

			if prefix == "" && pattern == "" && namespace == "" && dimensionName == "" && dimensionValue == "" {
				return fmt.Errorf("at least one filter must be specified")
			}

			return listMetricsWithFilters(cwClient, prefix, pattern, namespace, dimensionName, dimensionValue)
		},
	}

	listMetricsCmd.Flags().StringVarP(&metricName, "metric-name", "n", "", "Metric name to filter metrics")
	listMetricsCmd.Flags().StringVarP(&prefix, "prefix", "x", "", "Prefix to filter metrics by name")
	listMetricsCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter metrics by name")
	listMetricsCmd.Flags().StringVarP(&namespace, "namespace", "s", "", "Namespace to filter metrics")
	listMetricsCmd.Flags().StringVarP(&dimensionName, "dimension-name", "d", "", "Dimension name to filter metrics")
	listMetricsCmd.Flags().StringVarP(&dimensionValue, "dimension-value", "v", "", "Dimension value to filter metrics")
	listMetricsCmd.Flags().BoolVarP(&allMetrics, "all", "a", false, "List all metrics")

	cloudWatchCmd.AddCommand(listMetricsCmd)
}

func listAllMetrics(cwClient *cloudwatch.Client) error {
	result, err := cwClient.ListMetrics(context.TODO(), &cloudwatch.ListMetricsInput{})
	if err != nil {
		return fmt.Errorf("could not list metrics: %w", err)
	}

	for _, metric := range result.Metrics {
		printMetric(metric)
	}

	return nil
}

func listMetricsByName(cwClient *cloudwatch.Client, metricName string) error {
	result, err := cwClient.ListMetrics(context.TODO(), &cloudwatch.ListMetricsInput{
		MetricName: &metricName,
	})
	if err != nil {
		return fmt.Errorf("could not list metrics: %w", err)
	}

	for _, metric := range result.Metrics {
		printMetric(metric)
	}

	return nil
}

func listMetricsWithFilters(cwClient *cloudwatch.Client, prefix, pattern, namespace, dimensionName, dimensionValue string) error {
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
		printMetric(metric)
	}

	return nil
}

func printMetric(metric types.Metric) {
	fmt.Printf("Metric: %s\n", *metric.MetricName)
	for _, dimension := range metric.Dimensions {
		fmt.Printf("  Dimension: %s = %s\n", *dimension.Name, *dimension.Value)
	}
}
