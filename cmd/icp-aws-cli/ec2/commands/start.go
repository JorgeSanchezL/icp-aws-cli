package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
)

func InitStartCommands(ec2Client *ec2.Client, ec2Cmd *cobra.Command) {
	var pattern string
	var tagKey string
	var tagValue string

	var startInstanceCmd = &cobra.Command{
		Use:   "start [instance-id]",
		Short: "Starts an EC2 instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return startInstance(ec2Client, args[0])
		},
	}

	var startInstancesCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts EC2 instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			if pattern != "" {
				return startInstancesByPattern(ec2Client, pattern)
			} else if tagKey != "" && tagValue != "" {
				return startInstancesByTag(ec2Client, tagKey, tagValue)
			} else {
				return fmt.Errorf("either pattern or tag key and value must be specified")
			}
		},
	}

	startInstancesCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter instances")
	startInstancesCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter instances")
	startInstancesCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter instances")

	ec2Cmd.AddCommand(startInstanceCmd)
	ec2Cmd.AddCommand(startInstancesCmd)
}

func startInstance(ec2Client *ec2.Client, instanceID string) error {
	_, err := ec2Client.StartInstances(context.TODO(), &ec2.StartInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return fmt.Errorf("error starting instance %s: %w", instanceID, err)
	}

	fmt.Printf("Instance %s started successfully\n", instanceID)
	return nil
}

func startInstancesByPattern(ec2Client *ec2.Client, pattern string) error {
	instanceIDs, err := getInstanceIDsByPattern(ec2Client, pattern)
	if err != nil {
		return err
	}

	_, err = ec2Client.StartInstances(context.TODO(), &ec2.StartInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return fmt.Errorf("error starting instances: %w", err)
	}

	fmt.Printf("Instances %v started successfully\n", instanceIDs)
	return nil
}

func startInstancesByTag(ec2Client *ec2.Client, key, value string) error {
	instanceIDs, err := getInstanceIDsByTag(ec2Client, key, value)
	if err != nil {
		return err
	}

	_, err = ec2Client.StartInstances(context.TODO(), &ec2.StartInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return fmt.Errorf("error starting instances: %w", err)
	}

	fmt.Printf("Instances %v started successfully\n", instanceIDs)
	return nil
}
