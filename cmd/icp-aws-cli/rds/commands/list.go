package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/spf13/cobra"
)

func InitListCommands(rdsClient *rds.Client, rdsCmd *cobra.Command) {
	listInstancesCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists RDS instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listInstances(rdsClient)
		},
	}

	listSnapshotsCmd := &cobra.Command{
		Use:   "listSnapshots",
		Short: "Lists database snapshots",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listSnapshots(rdsClient, args[0])
		},
	}

	rdsCmd.AddCommand(listInstancesCmd)
	rdsCmd.AddCommand(listSnapshotsCmd)
}

func listInstances(rdsClient *rds.Client) error {
	result, err := rdsClient.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{})
	if err != nil {
		return fmt.Errorf("error listing RDS instances: %w", err)
	}

	for _, instance := range result.DBInstances {
		fmt.Println(*instance.DBInstanceIdentifier)
	}
	return nil
}

func listSnapshots(rdsClient *rds.Client, databaseID string) error {
	result, err := rdsClient.DescribeDBSnapshots(context.TODO(), &rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: &databaseID,
	})
	if err != nil {
		return fmt.Errorf("error listing snapshots: %w", err)
	}

	for _, snapshot := range result.DBSnapshots {
		fmt.Println(*snapshot.DBSnapshotIdentifier)
	}
	return nil
}
