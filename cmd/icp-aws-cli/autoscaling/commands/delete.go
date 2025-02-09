package commands

import (
	"context"
	"fmt"
	"icp-aws-cli/pkg/utils"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/spf13/cobra"
)

func InitDeleteGroupCommand(asClient *autoscaling.Client, autoscalingCmd *cobra.Command) {
	var groupName string
	var pattern string
	var tagKey string
	var tagValue string
	var allGroups bool

	var deleteGroupCmd = &cobra.Command{
		Use:   "delete-group",
		Short: "Deletes an AutoScaling group",
		RunE: func(cmd *cobra.Command, args []string) error {
			if allGroups && (groupName != "" || pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("the --all flag cannot be combined with other filters")
			}

			if allGroups {
				if !utils.ConfirmAction() {
					return fmt.Errorf("action cancelled by user")
				}
				return deleteAllGroups(asClient)
			}

			if groupName != "" && (pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("group name cannot be combined with other filters")
			}

			if groupName != "" {
				return deleteGroupByName(asClient, groupName)
			}

			if pattern == "" && tagKey == "" && tagValue == "" {
				return fmt.Errorf("at least one filter must be specified")
			}

			return deleteGroupsWithFilters(asClient, pattern, tagKey, tagValue)
		},
	}

	deleteGroupCmd.Flags().StringVarP(&groupName, "group-name", "g", "", "Group name to filter groups")
	deleteGroupCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter groups by name")
	deleteGroupCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter groups")
	deleteGroupCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter groups")
	deleteGroupCmd.Flags().BoolVarP(&allGroups, "all", "a", false, "Delete all groups")

	autoscalingCmd.AddCommand(deleteGroupCmd)
}

func deleteAllGroups(asClient *autoscaling.Client) error {
	result, err := asClient.DescribeAutoScalingGroups(context.TODO(), &autoscaling.DescribeAutoScalingGroupsInput{})
	if err != nil {
		return fmt.Errorf("could not list groups: %w", err)
	}

	for _, group := range result.AutoScalingGroups {
		_, err := asClient.DeleteAutoScalingGroup(context.TODO(), &autoscaling.DeleteAutoScalingGroupInput{
			AutoScalingGroupName: group.AutoScalingGroupName,
		})
		if err != nil {
			return fmt.Errorf("could not delete group %s: %w", *group.AutoScalingGroupName, err)
		}
		fmt.Printf("Deleted group %s\n", *group.AutoScalingGroupName)
	}

	return nil
}

func deleteGroupByName(asClient *autoscaling.Client, groupName string) error {
	_, err := asClient.DeleteAutoScalingGroup(context.TODO(), &autoscaling.DeleteAutoScalingGroupInput{
		AutoScalingGroupName: &groupName,
	})
	if err != nil {
		return fmt.Errorf("could not delete group %s: %w", groupName, err)
	}

	fmt.Printf("Deleted group %s\n", groupName)
	return nil
}

func deleteGroupsWithFilters(asClient *autoscaling.Client, pattern, tagKey, tagValue string) error {
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
		_, err := asClient.DeleteAutoScalingGroup(context.TODO(), &autoscaling.DeleteAutoScalingGroupInput{
			AutoScalingGroupName: group.AutoScalingGroupName,
		})
		if err != nil {
			return fmt.Errorf("could not delete group %s: %w", *group.AutoScalingGroupName, err)
		}
		fmt.Printf("Deleted group %s\n", *group.AutoScalingGroupName)
	}

	return nil
}
