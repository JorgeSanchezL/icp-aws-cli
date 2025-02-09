package metrics

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/spf13/cobra"
)

func InitCreateMetricCommand(cwClient *cloudwatch.Client, cloudWatchCmd *cobra.Command) {
	var metricName string
	var namespace string
	var dimensionName string
	var dimensionValue string

	var createMetricCmd = &cobra.Command{
		Use:   "create-metric",
		Short: "Creates a CloudWatch metric",
		RunE: func(cmd *cobra.Command, args []string) error {
			if metricName == "" || namespace == "" || dimensionName == "" || dimensionValue == "" {
				return fmt.Errorf("metric name, namespace, dimension name, and dimension value must be specified")
			}
			return createMetric(cwClient, metricName, namespace, dimensionName, dimensionValue)
		},
	}

	createMetricCmd.Flags().StringVarP(&metricName, "metric-name", "n", "", "Name of the metric to create")
	createMetricCmd.Flags().StringVarP(&namespace, "namespace", "s", "", "Namespace of the metric to create")
	createMetricCmd.Flags().StringVarP(&dimensionName, "dimension-name", "d", "", "Dimension name of the metric to create")
	createMetricCmd.Flags().StringVarP(&dimensionValue, "dimension-value", "v", "", "Dimension value of the metric to create")
	cloudWatchCmd.AddCommand(createMetricCmd)
}

func createMetric(cwClient *cloudwatch.Client, metricName, namespace, dimensionName, dimensionValue string) error {
	_, err := cwClient.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		Namespace: &namespace,
		MetricData: []types.MetricDatum{
			{
				MetricName: &metricName,
				Dimensions: []types.Dimension{
					{
						Name:  &dimensionName,
						Value: &dimensionValue,
					},
				},
				Value: aws.Float64(0), // Initial value
			},
		},
	})
	if err != nil {
		return fmt.Errorf("could not create metric %s: %w", metricName, err)
	}

	fmt.Printf("Created metric %s\n", metricName)
	return nil
}
