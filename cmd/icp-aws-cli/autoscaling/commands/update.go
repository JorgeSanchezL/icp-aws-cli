package commands

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/spf13/cobra"
)

func InitUpdateGroupCommand(asClient *autoscaling.Client, autoscalingCmd *cobra.Command) {
	var groupName string
	var pattern string
	var tagKey string
	var tagValue string
	var desiredCapacity int32

	var updateGroupCmd = &cobra.Command{
		Use:   "update-group",
		Short: "Updates an AutoScaling group",
		RunE: func(cmd *cobra.Command, args []string) error {
			if groupName != "" && (pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("group name cannot be combined with other filters")
			}

			if groupName != "" {
				return updateGroupByName(asClient, groupName, desiredCapacity)
			}

			if pattern == "" && tagKey == "" && tagValue == "" {
				return fmt.Errorf("at least one filter must be specified")
			}

			return updateGroupsWithFilters(asClient, pattern, tagKey, tagValue, desiredCapacity)
		},
	}

	updateGroupCmd.Flags().StringVarP(&groupName, "group-name", "g", "", "Group name to filter groups")
	updateGroupCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter groups by name")
	updateGroupCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter groups")
	updateGroupCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter groups")
	updateGroupCmd.Flags().Int32VarP(&desiredCapacity, "desired-capacity", "d", 0, "Desired capacity for the group")

	autoscalingCmd.AddCommand(updateGroupCmd)
}

func updateGroupByName(asClient *autoscaling.Client, groupName string, desiredCapacity int32) error {
	_, err := asClient.UpdateAutoScalingGroup(context.TODO(), &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: &groupName,
		DesiredCapacity:      &desiredCapacity,
	})
	if err != nil {
		return fmt.Errorf("could not update group %s: %w", groupName, err)
	}

	fmt.Printf("Updated group %s\n", groupName)
	return nil
}

func updateGroupsWithFilters(asClient *autoscaling.Client, pattern, tagKey, tagValue string, desiredCapacity int32) error {
	result, err := asClient.DescribeAutoScalingGroups(context.TODO(), &autoscaling.DescribeAutoScalingGroupsInput{})
	if err != nil {
		return fmt.Errorf("could not list groups: %w", err)
	}

	var groups []types.AutoScalingGroup
	if pattern != "" {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid pattern: %w", err)
		}
		for _, group := range result.AutoScalingGroups {
			if re.MatchString(*group.AutoScalingGroupName) {
				groups = append(groups, group)
			}
		}
	} else {
		groups = result.AutoScalingGroups
	}

	if tagKey != "" && tagValue != "" {
		var filteredGroups []types.AutoScalingGroup
		for _, group := range groups {
			for _, tag := range group.Tags {
				if *tag.Key == tagKey && *tag.Value == tagValue {
					filteredGroups = append(filteredGroups, group)
					break
				}
			}
		}
		groups = filteredGroups
	}

	for _, group := range groups {
		_, err := asClient.UpdateAutoScalingGroup(context.TODO(), &autoscaling.UpdateAutoScalingGroupInput{
			AutoScalingGroupName: group.AutoScalingGroupName,
			DesiredCapacity:      &desiredCapacity,
		})
		if err != nil {
			return fmt.Errorf("could not update group %s: %w", *group.AutoScalingGroupName, err)
		}
		fmt.Printf("Updated group %s\n", *group.AutoScalingGroupName)
	}

	return nil
}
