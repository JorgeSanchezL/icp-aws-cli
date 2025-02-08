package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
)

func InitListCommands(ec2Client *ec2.Client, ec2Cmd *cobra.Command) {
	var pattern string
	var tagKey string
	var tagValue string

	var listInstancesCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists EC2 instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			if pattern != "" {
				return listInstancesByPattern(ec2Client, pattern)
			} else if tagKey != "" && tagValue != "" {
				return listInstancesByTag(ec2Client, tagKey, tagValue)
			} else {
				return listInstances(ec2Client)
			}
		},
	}

	listInstancesCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter instances")
	listInstancesCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter instances")
	listInstancesCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter instances")

	var listRunningInstancesCmd = &cobra.Command{
		Use:   "list-running",
		Short: "Lists running EC2 instances with uptime",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRunningInstances(ec2Client)
		},
	}

	ec2Cmd.AddCommand(listInstancesCmd)
	ec2Cmd.AddCommand(listRunningInstancesCmd)
}

func listInstances(ec2Client *ec2.Client) error {
	result, err := ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return fmt.Errorf("error listing EC2 instances: %w", err)
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Println(*instance.InstanceId)
		}
	}
	return nil
}

func listRunningInstances(ec2Client *ec2.Client) error {
	result, err := ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error listing running instances: %w", err)
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			uptime := time.Since(*instance.LaunchTime).Hours()
			fmt.Printf("Instance ID: %s, Uptime: %.2f hours\n", *instance.InstanceId, uptime)
		}
	}
	return nil
}

func listInstancesByPattern(ec2Client *ec2.Client, pattern string) error {
	instanceIDs, err := getInstanceIDsByPattern(ec2Client, pattern)
	if err != nil {
		return err
	}

	for _, instanceID := range instanceIDs {
		fmt.Println(instanceID)
	}
	return nil
}

func listInstancesByTag(ec2Client *ec2.Client, key, value string) error {
	instanceIDs, err := getInstanceIDsByTag(ec2Client, key, value)
	if err != nil {
		return err
	}

	for _, instanceID := range instanceIDs {
		fmt.Println(instanceID)
	}
	return nil
}
