package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/spf13/cobra"
)

func InitUpdateGroupCommand(asClient *autoscaling.Client, autoscalingCmd *cobra.Command) {
	var groupName string
	var minSize int32
	var maxSize int32
	var desiredCapacity int32

	var updateGroupCmd = &cobra.Command{
		Use:   "update-group",
		Short: "Updates an AutoScaling group",
		RunE: func(cmd *cobra.Command, args []string) error {
			if groupName == "" {
				return fmt.Errorf("group name must be specified")
			}

			return updateGroup(asClient, groupName, minSize, maxSize, desiredCapacity)
		},
	}

	updateGroupCmd.Flags().StringVarP(&groupName, "group-name", "g", "", "AutoScaling group name to update")
	updateGroupCmd.Flags().Int32VarP(&minSize, "min-size", "m", 0, "Minimum size of the group")
	updateGroupCmd.Flags().Int32VarP(&maxSize, "max-size", "x", 0, "Maximum size of the group")
	updateGroupCmd.Flags().Int32VarP(&desiredCapacity, "desired-capacity", "d", 0, "Desired capacity of the group")

	autoscalingCmd.AddCommand(updateGroupCmd)
}

func updateGroup(asClient *autoscaling.Client, groupName string, minSize, maxSize, desiredCapacity int32) error {
	_, err := asClient.UpdateAutoScalingGroup(context.TODO(), &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(groupName),
		MinSize:              aws.Int32(minSize),
		MaxSize:              aws.Int32(maxSize),
		DesiredCapacity:      aws.Int32(desiredCapacity),
	})
	if err != nil {
		return fmt.Errorf("could not update AutoScaling group %s: %w", groupName, err)
	}

	fmt.Printf("Updated AutoScaling group %s\n", groupName)
	return nil
}
