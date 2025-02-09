package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/spf13/cobra"
)

func InitDeleteGroupsCommand(asClient *autoscaling.Client, autoscalingCmd *cobra.Command) {
	var groupName string
	var pattern string
	var tagKey string
	var tagValue string
	var allGroups bool

	var deleteGroupsCmd = &cobra.Command{
		Use:   "delete",
		Short: "Deletes AutoScaling groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			if allGroups && (groupName != "" || pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("the --all flag cannot be combined with other filters")
			}

			if allGroups {
				return deleteAllGroups(asClient)
			}

			if groupName != "" && (pattern != "" || tagKey != "" || tagValue != "") {
				return fmt.Errorf("group name cannot be combined with other filters")
			}

			if groupName != "" {
				return deleteGroupsByName(asClient, groupName)
			}

			if pattern == "" && tagKey == "" && tagValue == "" {
				return fmt.Errorf("at least one filter must be specified")
			}

			return deleteGroupsWithFilters(asClient, pattern, tagKey, tagValue)
		},
	}

	deleteGroupsCmd.Flags().StringVarP(&groupName, "group-name", "g", "", "AutoScaling group name to filter groups")
	deleteGroupsCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter groups by name")
	deleteGroupsCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter groups")
	deleteGroupsCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter groups")
	deleteGroupsCmd.Flags().BoolVarP(&allGroups, "all", "a", false, "Delete all groups")

	autoscalingCmd.AddCommand(deleteGroupsCmd)
}

func deleteAllGroups(asClient *autoscaling.Client) error {
	result, err := asClient.DescribeAutoScalingGroups(context.TODO(), &autoscaling.DescribeAutoScalingGroupsInput{})
	if err != nil {
		return fmt.Errorf("could not list AutoScaling groups: %w", err)
	}

	for _, group := range result.AutoScalingGroups {
		if err := deleteGroup(asClient, *group.AutoScalingGroupName); err != nil {
			return err
		}
	}

	return nil
}

func deleteGroupsByName(asClient *autoscaling.Client, groupName string) error {
	return deleteGroup(asClient, groupName)
}

func deleteGroupsWithFilters(asClient *autoscaling.Client, pattern, tagKey, tagValue string) error {
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
		if err := deleteGroup(asClient, *group.AutoScalingGroupName); err != nil {
			return err
		}
	}

	return nil
}

func deleteGroup(asClient *autoscaling.Client, groupName string) error {
	_, err := asClient.DeleteAutoScalingGroup(context.TODO(), &autoscaling.DeleteAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(groupName),
		ForceDelete:          aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("could not delete AutoScaling group %s: %w", groupName, err)
	}

	fmt.Printf("Deleted AutoScaling group %s\n", groupName)
	return nil
}
