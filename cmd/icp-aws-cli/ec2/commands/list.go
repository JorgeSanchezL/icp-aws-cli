package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
)

func InitListCommands(ec2Client *ec2.Client, ec2Cmd *cobra.Command) {
	var instanceID string
	var pattern string
	var tagKey string
	var tagValue string

	var listInstancesCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists EC2 instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			if instanceID != "" && (pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("instance ID cannot be combined with other filters")
			}

			filters := []types.Filter{}

			if instanceID != "" {
				filters = append(filters, types.Filter{
					Name:   aws.String("instance-id"),
					Values: []string{instanceID},
				})
			}

			if pattern != "" {
				filters = append(filters, types.Filter{
					Name:   aws.String("tag:Name"),
					Values: []string{pattern},
				})
			}

			if tagKey != "" {
				if tagValue == "" {
					return fmt.Errorf("tag value must be specified when tag key is provided")
				}
				filters = append(filters, types.Filter{
					Name:   aws.String(fmt.Sprintf("tag:%s", tagKey)),
					Values: []string{tagValue},
				})
			}

			return manageInstancesWithFilters(ec2Client, filters, buildListInstancesInput, listInstances)
		},
	}

	listInstancesCmd.Flags().StringVarP(&instanceID, "instance-id", "i", "", "Instance ID to filter instances")
	listInstancesCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter instances")
	listInstancesCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter instances")
	listInstancesCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter instances")

	ec2Cmd.AddCommand(listInstancesCmd)
}

func buildListInstancesInput(instanceIDs []string) interface{} {
	return &ec2.DescribeInstancesInput{
		InstanceIds: instanceIDs,
	}
}

func listInstances(ec2Client *ec2.Client, ctx context.Context, input interface{}) (interface{}, error) {
	ec2Input := input.(*ec2.DescribeInstancesInput)
	return ec2Client.DescribeInstances(ctx, ec2Input)
}
