package commands

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/spf13/cobra"
)

func InitListGroupsCommand(asClient *autoscaling.Client, autoscalingCmd *cobra.Command) {
	var groupName string
	var pattern string
	var tagKey string
	var tagValue string
	var allGroups bool

	var listGroupsCmd = &cobra.Command{
		Use:   "list-groups",
		Short: "Lists AutoScaling groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			if allGroups && (groupName != "" || pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("the --all flag cannot be combined with other filters")
			}

			if allGroups {
				return listAllGroups(asClient)
			}

			if groupName != "" && (pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("group name cannot be combined with other filters")
			}

			if groupName != "" {
				return listGroupsByName(asClient, groupName)
			}

			if pattern == "" && tagKey == "" && tagValue == "" {
				return fmt.Errorf("at least one filter must be specified")
			}

			return listGroupsWithFilters(asClient, pattern, tagKey, tagValue)
		},
	}

	listGroupsCmd.Flags().StringVarP(&groupName, "group-name", "g", "", "AutoScaling group name to filter groups")
	listGroupsCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter groups by name")
	listGroupsCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter groups")
	listGroupsCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter groups")
	listGroupsCmd.Flags().BoolVarP(&allGroups, "all", "a", false, "List all groups")

	autoscalingCmd.AddCommand(listGroupsCmd)
}

func listAllGroups(asClient *autoscaling.Client) error {
	result, err := asClient.DescribeAutoScalingGroups(context.TODO(), &autoscaling.DescribeAutoScalingGroupsInput{})
	if err != nil {
		return fmt.Errorf("could not list AutoScaling groups: %w", err)
	}

	for _, group := range result.AutoScalingGroups {
		printGroup(group)
	}

	return nil
}

func listGroupsByName(asClient *autoscaling.Client, groupName string) error {
	result, err := asClient.DescribeAutoScalingGroups(context.TODO(), &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{groupName},
	})
	if err != nil {
		return fmt.Errorf("could not list AutoScaling groups: %w", err)
	}

	for _, group := range result.AutoScalingGroups {
		printGroup(group)
	}

	return nil
}

func listGroupsWithFilters(asClient *autoscaling.Client, pattern, tagKey, tagValue string) error {
	result, err := asClient.DescribeAutoScalingGroups(context.TODO(), &autoscaling.DescribeAutoScalingGroupsInput{})
	if err != nil {
		return fmt.Errorf("could not list AutoScaling groups: %w", err)
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
		printGroup(group)
	}

	return nil
}

func printGroup(group types.AutoScalingGroup) {
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
