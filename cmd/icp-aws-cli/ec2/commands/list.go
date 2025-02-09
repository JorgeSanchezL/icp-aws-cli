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
	var allInstances bool
	var state string

	var listInstancesCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists EC2 instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			if allInstances && (instanceID != "" || pattern != "" || tagKey != "" || tagValue != "" || state != "") {
				return fmt.Errorf("the --all flag cannot be combined with other filters")
			}

			if allInstances {
				return manageInstancesWithFilters(ec2Client, []types.Filter{}, buildListInstancesInput, listInstances)
			}

			if instanceID != "" && (pattern != "" || tagKey != "" || tagValue != "" || state != "") {
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

			if state != "" {
				filters = append(filters, types.Filter{
					Name:   aws.String("instance-state-name"),
					Values: []string{state},
				})
			}

			if len(filters) == 0 {
				return fmt.Errorf("at least one filter must be specified")
			}

			return manageInstancesWithFilters(ec2Client, filters, buildListInstancesInput, listInstances)
		},
	}

	listInstancesCmd.Flags().StringVarP(&instanceID, "instance-id", "i", "", "Instance ID to filter instances")
	listInstancesCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter instances")
	listInstancesCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter instances")
	listInstancesCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter instances")
	listInstancesCmd.Flags().BoolVarP(&allInstances, "all", "a", false, "Apply action to all instances")
	listInstancesCmd.Flags().StringVarP(&state, "state", "s", "", "State to filter instances (e.g., running, stopped)")
	ec2Cmd.AddCommand(listInstancesCmd)
}

func buildListInstancesInput(instanceIDs []string) interface{} {
	return &ec2.DescribeInstancesInput{
		InstanceIds: instanceIDs,
	}
}

func listInstances(ec2Client *ec2.Client, ctx context.Context, input interface{}) (interface{}, error) {
	ec2Input := input.(*ec2.DescribeInstancesInput)
	output, err := ec2Client.DescribeInstances(ctx, ec2Input)
	if err != nil {
		return nil, fmt.Errorf("error describing instances: %w", err)
	}

	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			name := "<Not Assigned>"
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}
			fmt.Printf("Name: %s, ID: %s, Type: %s, State: %s, Launched: %s\n", name, *instance.InstanceId, instance.InstanceType, instance.State.Name, instance.LaunchTime.Format("2006-01-02 15:04:05"))
		}
	}

	return nil, nil
}
