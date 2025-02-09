package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/spf13/cobra"
)

func InitCreateGroupCommand(asClient *autoscaling.Client, autoscalingCmd *cobra.Command) {
	var groupName, launchConfigurationName string
	var minSize, maxSize, desiredCapacity int32
	var tags []string

	var createGroupCmd = &cobra.Command{
		Use:   "create",
		Short: "Creates an AutoScaling group",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createGroup(asClient, groupName, launchConfigurationName, minSize, maxSize, desiredCapacity, tags)
		},
	}

	createGroupCmd.Flags().StringVarP(&groupName, "group-name", "g", "", "Name of the AutoScaling group")
	createGroupCmd.Flags().StringVarP(&launchConfigurationName, "launch-configuration-name", "l", "", "Launch configuration name")
	createGroupCmd.Flags().Int32VarP(&minSize, "min-size", "m", 1, "Minimum size of the group")
	createGroupCmd.Flags().Int32VarP(&maxSize, "max-size", "x", 1, "Maximum size of the group")
	createGroupCmd.Flags().Int32VarP(&desiredCapacity, "desired-capacity", "d", 1, "Desired capacity of the group")
	createGroupCmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, "Tags for the group (key=value)")

	autoscalingCmd.AddCommand(createGroupCmd)
}

func createGroup(asClient *autoscaling.Client, groupName, launchConfigurationName string, minSize, maxSize, desiredCapacity int32, tags []string) error {
	tagList := []types.Tag{}
	for _, tag := range tags {
		parts := strings.SplitN(tag, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid tag format: %s", tag)
		}
		tagList = append(tagList, types.Tag{
			Key:   &parts[0],
			Value: &parts[1],
		})
	}

	input := &autoscaling.CreateAutoScalingGroupInput{
		AutoScalingGroupName:    &groupName,
		LaunchConfigurationName: &launchConfigurationName,
		MinSize:                 &minSize,
		MaxSize:                 &maxSize,
		DesiredCapacity:         &desiredCapacity,
		Tags:                    tagList,
	}

	_, err := asClient.CreateAutoScalingGroup(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("could not create AutoScaling group: %w", err)
	}

	fmt.Printf("Created AutoScaling group %s\n", groupName)
	return nil
}
