package rds

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/spf13/cobra"
)

func InitCommands(rdsClient *rds.Client) *cobra.Command {
	var rdsCmd = &cobra.Command{
		Use:   "rds",
		Short: "Commands to interact with Amazon RDS",
		Long:  "Allows listing and managing RDS instances in Amazon RDS.",
	}

	var listInstancesCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists RDS instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listInstances(rdsClient)
		},
	}

	rdsCmd.AddCommand(listInstancesCmd)
	return rdsCmd
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
