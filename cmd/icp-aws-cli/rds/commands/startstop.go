package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/spf13/cobra"
)

func InitStartStopCommands(rdsClient *rds.Client, rdsCmd *cobra.Command) {
	startInstanceCmd := &cobra.Command{
		Use:   "startInstance",
		Short: "Starts a stopped RDS instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return startInstance(rdsClient, args[0])
		},
	}

	stopInstanceCmd := &cobra.Command{
		Use:   "stopInstance",
		Short: "Stops a running RDS instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return stopInstance(rdsClient, args[0])
		},
	}

	rdsCmd.AddCommand(startInstanceCmd)
	rdsCmd.AddCommand(stopInstanceCmd)
}

func startInstance(rdsClient *rds.Client, instanceID string) error {
	_, err := rdsClient.StartDBInstance(context.TODO(), &rds.StartDBInstanceInput{
		DBInstanceIdentifier: &instanceID,
	})
	if err != nil {
		return fmt.Errorf("error starting instance: %w", err)
	}

	fmt.Printf("Instance %s starting\n", instanceID)
	return nil
}

func stopInstance(rdsClient *rds.Client, instanceID string) error {
	_, err := rdsClient.StopDBInstance(context.TODO(), &rds.StopDBInstanceInput{
		DBInstanceIdentifier: &instanceID,
	})
	if err != nil {
		return fmt.Errorf("error stopping instance: %w", err)
	}

	fmt.Printf("Instance %s stopping\n", instanceID)
	return nil
}
