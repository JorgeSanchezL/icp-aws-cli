package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
)

func InitTerminateCommands(ec2Client *ec2.Client, ec2Cmd *cobra.Command) {
	var pattern string
	var tagKey string
	var tagValue string

	var terminateInstanceCmd = &cobra.Command{
		Use:   "terminate [instance-id]",
		Short: "Terminates an EC2 instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return terminateInstance(ec2Client, args[0])
		},
	}

	var terminateInstancesCmd = &cobra.Command{
		Use:   "terminate",
		Short: "Terminates EC2 instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			if pattern != "" {
				return terminateInstancesByPattern(ec2Client, pattern)
			} else if tagKey != "" && tagValue != "" {
				return terminateInstancesByTag(ec2Client, tagKey, tagValue)
			} else {
				return fmt.Errorf("either pattern or tag key and value must be specified")
			}
		},
	}

	terminateInstancesCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter instances")
	terminateInstancesCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter instances")
	terminateInstancesCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter instances")

	ec2Cmd.AddCommand(terminateInstanceCmd)
	ec2Cmd.AddCommand(terminateInstancesCmd)
}

func terminateInstance(ec2Client *ec2.Client, instanceID string) error {
	_, err := ec2Client.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return fmt.Errorf("error terminating instance %s: %w", instanceID, err)
	}

	fmt.Printf("Instance %s terminated successfully\n", instanceID)
	return nil
}

func terminateInstancesByPattern(ec2Client *ec2.Client, pattern string) error {
	instanceIDs, err := getInstanceIDsByPattern(ec2Client, pattern)
	if err != nil {
		return err
	}

	_, err = ec2Client.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return fmt.Errorf("error terminating instances: %w", err)
	}

	fmt.Printf("Instances %v terminated successfully\n", instanceIDs)
	return nil
}

func terminateInstancesByTag(ec2Client *ec2.Client, key, value string) error {
	instanceIDs, err := getInstanceIDsByTag(ec2Client, key, value)
	if err != nil {
		return err
	}

	_, err = ec2Client.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return fmt.Errorf("error terminating instances: %w", err)
	}

	fmt.Printf("Instances %v terminated successfully\n", instanceIDs)
	return nil
}
