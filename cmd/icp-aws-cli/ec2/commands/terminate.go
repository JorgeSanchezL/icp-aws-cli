package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
)

func InitTerminateCommands(ec2Client *ec2.Client, ec2Cmd *cobra.Command) {
	var instanceID string
	var pattern string
	var tagKey string
	var tagValue string
	var allInstances bool

	var terminateInstancesCmd = &cobra.Command{
		Use:   "terminate",
		Short: "Terminates EC2 instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			if allInstances && (instanceID != "" || pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("the --all flag cannot be combined with other filters")
			}

			if allInstances {
				if !confirmAction() {
					return fmt.Errorf("action cancelled by user")
				}
				return manageInstancesWithFilters(ec2Client, []types.Filter{}, buildTerminateInstancesInput, terminateInstances)
			}

			if instanceID != "" && (pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("instance ID cannot be combined with other filters")
			}

			if instanceID != "" {
				return terminateInstancesByID(ec2Client, instanceID)
			}

			filters := []types.Filter{}

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

			if len(filters) == 0 {
				return fmt.Errorf("at least one filter must be specified")
			}

			return manageInstancesWithFilters(ec2Client, filters, buildTerminateInstancesInput, terminateInstances)
		},
	}

	terminateInstancesCmd.Flags().StringVarP(&instanceID, "instance-id", "i", "", "Instance ID to filter instances")
	terminateInstancesCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter instances")
	terminateInstancesCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter instances")
	terminateInstancesCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter instances")
	terminateInstancesCmd.Flags().BoolVarP(&allInstances, "all", "a", false, "Apply action to all instances")

	ec2Cmd.AddCommand(terminateInstancesCmd)
}

func terminateInstancesByID(ec2Client *ec2.Client, instanceID string) error {
	_, err := ec2Client.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return fmt.Errorf("error terminating instance %s: %w", instanceID, err)
	}
	fmt.Printf("Instance %s terminated successfully\n", instanceID)
	return nil
}

func buildTerminateInstancesInput(instanceIDs []string) interface{} {
	return &ec2.TerminateInstancesInput{
		InstanceIds: instanceIDs,
	}
}

func terminateInstances(ec2Client *ec2.Client, ctx context.Context, input interface{}) (interface{}, error) {
	ec2Input := input.(*ec2.TerminateInstancesInput)
	return ec2Client.TerminateInstances(ctx, ec2Input)
}
