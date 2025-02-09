package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/spf13/cobra"
)

func InitDescribeGroupCommand(asClient *autoscaling.Client, autoscalingCmd *cobra.Command) {
	var groupName string

	var describeGroupCmd = &cobra.Command{
		Use:   "describe-group",
		Short: "Describes an AutoScaling group",
		RunE: func(cmd *cobra.Command, args []string) error {
			return describeGroup(asClient, groupName)
		},
	}

	describeGroupCmd.Flags().StringVarP(&groupName, "group-name", "g", "", "Name of the AutoScaling group")
	autoscalingCmd.AddCommand(describeGroupCmd)
}

func describeGroup(asClient *autoscaling.Client, groupName string) error {
	if groupName == "" {
		return fmt.Errorf("group name must be specified")
	}

	result, err := asClient.DescribeAutoScalingGroups(context.TODO(), &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{groupName},
	})
	if err != nil {
		return fmt.Errorf("could not describe AutoScaling group: %w", err)
	}

	for _, group := range result.AutoScalingGroups {
		fmt.Printf("Group: %s\n", *group.AutoScalingGroupName)
		fmt.Printf("Launch Configuration: %s\n", *group.LaunchConfigurationName)
		fmt.Printf("Min Size: %d\n", group.MinSize)
		fmt.Printf("Max Size: %d\n", group.MaxSize)
		fmt.Printf("Desired Capacity: %d\n", group.DesiredCapacity)
		fmt.Printf("Instances: %d\n", len(group.Instances))
		for _, tag := range group.Tags {
			fmt.Printf("  Tag: %s = %s\n", *tag.Key, *tag.Value)
		}
	}

	return nil
}
