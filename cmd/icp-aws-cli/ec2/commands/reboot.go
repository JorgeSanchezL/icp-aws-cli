package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
)

func InitRebootCommands(ec2Client *ec2.Client, ec2Cmd *cobra.Command) {
	var pattern string
	var tagKey string
	var tagValue string

	var rebootInstanceCmd = &cobra.Command{
		Use:   "reboot [instance-id]",
		Short: "Reboots an EC2 instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return rebootInstance(ec2Client, args[0])
		},
	}

	var rebootInstancesCmd = &cobra.Command{
		Use:   "reboot",
		Short: "Reboots EC2 instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			if pattern != "" {
				return rebootInstancesByPattern(ec2Client, pattern)
			} else if tagKey != "" && tagValue != "" {
				return rebootInstancesByTag(ec2Client, tagKey, tagValue)
			} else {
				return fmt.Errorf("either pattern or tag key and value must be specified")
			}
		},
	}

	rebootInstancesCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern to filter instances")
	rebootInstancesCmd.Flags().StringVarP(&tagKey, "tag-key", "k", "", "Tag key to filter instances")
	rebootInstancesCmd.Flags().StringVarP(&tagValue, "tag-value", "v", "", "Tag value to filter instances")

	ec2Cmd.AddCommand(rebootInstanceCmd)
	ec2Cmd.AddCommand(rebootInstancesCmd)
}

func rebootInstance(ec2Client *ec2.Client, instanceID string) error {
	_, err := ec2Client.RebootInstances(context.TODO(), &ec2.RebootInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return fmt.Errorf("error rebooting instance %s: %w", instanceID, err)
	}

	fmt.Printf("Instance %s rebooted successfully\n", instanceID)
	return nil
}

func rebootInstancesByPattern(ec2Client *ec2.Client, pattern string) error {
	instanceIDs, err := getInstanceIDsByPattern(ec2Client, pattern)
	if err != nil {
		return err
	}

	_, err = ec2Client.RebootInstances(context.TODO(), &ec2.RebootInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return fmt.Errorf("error rebooting instances: %w", err)
	}

	fmt.Printf("Instances %v rebooted successfully\n", instanceIDs)
	return nil
}

func rebootInstancesByTag(ec2Client *ec2.Client, key, value string) error {
	instanceIDs, err := getInstanceIDsByTag(ec2Client, key, value)
	if err != nil {
		return err
	}

	_, err = ec2Client.RebootInstances(context.TODO(), &ec2.RebootInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return fmt.Errorf("error rebooting instances: %w", err)
	}

	fmt.Printf("Instances %v rebooted successfully\n", instanceIDs)
	return nil
}
