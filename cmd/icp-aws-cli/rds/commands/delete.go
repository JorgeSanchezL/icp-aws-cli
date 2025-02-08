package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/spf13/cobra"
)

func InitDeleteCommands(rdsClient *rds.Client, rdsCmd *cobra.Command) {
	deleteInstanceCmd := &cobra.Command{
		Use:   "deleteInstance",
		Short: "Deletes an RDS instance",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteInstance(rdsClient, args[0], args[1])
		},
	}

	deleteSnapshotCmd := &cobra.Command{
		Use:   "deleteSnapshot",
		Short: "Deletes a database snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteSnapshot(rdsClient, args[0])
		},
	}

	rdsCmd.AddCommand(deleteInstanceCmd)
	rdsCmd.AddCommand(deleteSnapshotCmd)
}

func deleteInstance(rdsClient *rds.Client, databaseName string, skipFinalSnapshot string) error {
	skip, err := strconv.ParseBool(skipFinalSnapshot)
	if err != nil {
		return fmt.Errorf("invalid skip snapshot value: %w", err)
	}

	_, err = rdsClient.DeleteDBInstance(context.TODO(), &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: &databaseName,
		SkipFinalSnapshot:    &skip,
	})
	if err != nil {
		return fmt.Errorf("error deleting instance: %w", err)
	}

	fmt.Printf("Instance %s deletion initiated\n", databaseName)
	return nil
}

func deleteSnapshot(rdsClient *rds.Client, snapshotID string) error {
	_, err := rdsClient.DeleteDBSnapshot(context.TODO(), &rds.DeleteDBSnapshotInput{
		DBSnapshotIdentifier: &snapshotID,
	})
	if err != nil {
		return fmt.Errorf("error deleting snapshot: %w", err)
	}

	fmt.Printf("Snapshot %s deleted\n", snapshotID)
	return nil
}
