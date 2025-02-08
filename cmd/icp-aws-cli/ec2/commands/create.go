package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
)

func InitCreateCommands(ec2Client *ec2.Client, ec2Cmd *cobra.Command) {
	var createInstanceCmd = &cobra.Command{
		Use:   "create [ami-id] [instance-type]",
		Short: "Creates a new EC2 instance",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createInstance(ec2Client, args[0], args[1])
		},
	}

	ec2Cmd.AddCommand(createInstanceCmd)
}

func createInstance(ec2Client *ec2.Client, amiID, instanceType string) error {
	runResult, err := ec2Client.RunInstances(context.TODO(), &ec2.RunInstancesInput{
		ImageId:      aws.String(amiID),
		InstanceType: types.InstanceType(instanceType),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
	})
	if err != nil {
		return fmt.Errorf("could not create instance: %w", err)
	}

	fmt.Printf("Created instance %s\n", *runResult.Instances[0].InstanceId)
	return nil
}
