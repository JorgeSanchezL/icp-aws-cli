package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/spf13/cobra"
)

func InitCreateCommands(rdsClient *rds.Client, rdsCmd *cobra.Command) {
	createSnapshotCmd := &cobra.Command{
		Use:   "createSnapshot",
		Short: "Creates a database snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createSnapshot(rdsClient, args[0])
		},
	}

	createInstanceCmd := &cobra.Command{
		Use:   "createInstance",
		Short: "Creates a new RDS instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createInstance(rdsClient, args[0])
		},
	}

	rdsCmd.AddCommand(createSnapshotCmd)
	rdsCmd.AddCommand(createInstanceCmd)
}

func createSnapshot(rdsClient *rds.Client, configJSON string) error {
	var input rds.CreateDBSnapshotInput

	if err := json.Unmarshal([]byte(configJSON), &input); err != nil {
		return fmt.Errorf("error parsing input JSON: %w", err)
	}

	_, err := rdsClient.CreateDBSnapshot(context.TODO(), &input)
	if err != nil {
		return fmt.Errorf("error creating snapshot: %w", err)
	}

	fmt.Printf("Snapshot %s created\n", *input.DBSnapshotIdentifier)
	return nil
}

func createInstance(rdsClient *rds.Client, configJSON string) error {
	var input rds.CreateDBInstanceInput

	if err := json.Unmarshal([]byte(configJSON), &input); err != nil {
		return fmt.Errorf("error parsing input JSON: %w", err)
	}

	_, err := rdsClient.CreateDBInstance(context.TODO(), &input)
	if err != nil {
		return fmt.Errorf("error creating instance: %w", err)
	}

	fmt.Printf("Instance %s creation started\n", *input.DBInstanceIdentifier)
	return nil
}
