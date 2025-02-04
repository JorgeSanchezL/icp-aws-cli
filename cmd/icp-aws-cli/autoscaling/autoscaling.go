package autoscaling

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/spf13/cobra"
)

func InitCommands(autoscalingClient *autoscaling.Client) *cobra.Command {
	var autoscalingCmd = &cobra.Command{
		Use:   "autoscaling",
		Short: "Commands to interact with Amazon AutoScaling",
		Long:  "Allows listing and managing AutoScaling groups in Amazon AutoScaling.",
	}

	var listGroupsCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists AutoScaling groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listGroups(autoscalingClient)
		},
	}

	autoscalingCmd.AddCommand(listGroupsCmd)
	return autoscalingCmd
}

func listGroups(autoscalingClient *autoscaling.Client) error {
	result, err := autoscalingClient.DescribeAutoScalingGroups(context.TODO(), &autoscaling.DescribeAutoScalingGroupsInput{})
	if err != nil {
		return fmt.Errorf("error listing AutoScaling groups: %w", err)
	}

	for _, group := range result.AutoScalingGroups {
		fmt.Println(*group.AutoScalingGroupName)
	}
	return nil
}
