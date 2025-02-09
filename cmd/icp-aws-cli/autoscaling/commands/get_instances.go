package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/spf13/cobra"
)

func InitGetInstancesCommand(asClient *autoscaling.Client, autoscalingCmd *cobra.Command) {
	var groupName string

	var getInstancesCmd = &cobra.Command{
		Use:   "get-instances",
		Short: "Gets the instance IDs of all instances in an AutoScaling group",
		RunE: func(cmd *cobra.Command, args []string) error {
			if groupName == "" {
				return fmt.Errorf("group name must be specified")
			}
			return getInstances(asClient, groupName)
		},
	}

	getInstancesCmd.Flags().StringVarP(&groupName, "group-name", "g", "", "AutoScaling group name")
	autoscalingCmd.AddCommand(getInstancesCmd)
}

func getInstances(asClient *autoscaling.Client, groupName string) error {
	result, err := asClient.DescribeAutoScalingGroups(context.TODO(), &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{groupName},
	})
	if err != nil {
		return fmt.Errorf("could not describe AutoScaling group: %w", err)
	}

	if len(result.AutoScalingGroups) == 0 {
		return fmt.Errorf("no AutoScaling group found with name %s", groupName)
	}

	group := result.AutoScalingGroups[0]
	for _, instance := range group.Instances {
		fmt.Printf("Instance ID: %s\n", *instance.InstanceId)
	}

	return nil
}
