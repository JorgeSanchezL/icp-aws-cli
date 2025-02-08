package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
)

func InitStopCommands(ec2Client *ec2.Client, ec2Cmd *cobra.Command) {
	var pattern string
	var tagKey string
	var tagValue string

	var stopInstanceCmd = &cobra.Command{
		Use:   "stop [instance-id]",
		Short: "Stops an EC2 instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return stopInstance(ec2Client, args[0])
		},
	}

	var stopInstancesCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stops EC2 instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			if pattern != "" {
				return stopInstancesByPattern(ec2Client, pattern)
			} else if tagKey != "" && tagValue != "" {
				return stopInstancesByTag(ec2Client, tagKey, tagValue)
			} else {
				return fmt.Errorf("either pattern or tag key and value must be specified")
			}
		},
	}

	stopInstancesCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter instances")
	stopInstancesCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter instances")
	stopInstancesCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter instances")

	ec2Cmd.AddCommand(stopInstanceCmd)
	ec2Cmd.AddCommand(stopInstancesCmd)
}

func stopInstance(ec2Client *ec2.Client, instanceID string) error {
	_, err := ec2Client.StopInstances(context.TODO(), &ec2.StopInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return fmt.Errorf("error stopping instance %s: %w", instanceID, err)
	}

	fmt.Printf("Instance %s stopped successfully\n", instanceID)
	return nil
}

func stopInstancesByPattern(ec2Client *ec2.Client, pattern string) error {
	instanceIDs, err := getInstanceIDsByPattern(ec2Client, pattern)
	if err != nil {
		return err
	}

	_, err = ec2Client.StopInstances(context.TODO(), &ec2.StopInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return fmt.Errorf("error stopping instances: %w", err)
	}

	fmt.Printf("Instances %v stopped successfully\n", instanceIDs)
	return nil
}

func stopInstancesByTag(ec2Client *ec2.Client, key, value string) error {
	instanceIDs, err := getInstanceIDsByTag(ec2Client, key, value)
	if err != nil {
		return err
	}

	_, err = ec2Client.StopInstances(context.TODO(), &ec2.StopInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return fmt.Errorf("error stopping instances: %w", err)
	}

	fmt.Printf("Instances %v stopped successfully\n", instanceIDs)
	return nil
}
