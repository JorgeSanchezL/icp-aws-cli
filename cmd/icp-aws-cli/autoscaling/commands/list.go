package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
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

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		Filters: filters,
	}

	result, err := asClient.DescribeAutoScalingGroups(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("could not list AutoScaling groups: %w", err)
	}

	for _, group := range result.AutoScalingGroups {
		printGroup(group)
	}

	return nil
}

func printGroup(group types.AutoScalingGroup) {
	fmt.Printf("Group: %s, Current number of instances: %d, MinSize: %d, MaxSize: %d, DesiredCapacity: %d\n", *group.AutoScalingGroupName, len(group.Instances), group.MinSize, group.MaxSize, group.DesiredCapacity)
}
