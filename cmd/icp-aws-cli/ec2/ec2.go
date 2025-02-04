package ec2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
)

func InitCommands(ec2Client *ec2.Client) *cobra.Command {
	var ec2Cmd = &cobra.Command{
		Use:   "ec2",
		Short: "Commands to interact with Amazon EC2",
		Long:  "Allows listing and managing EC2 instances in Amazon EC2.",
	}

	var listInstancesCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists EC2 instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listInstances(ec2Client)
		},
	}

	ec2Cmd.AddCommand(listInstancesCmd)
	return ec2Cmd
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
