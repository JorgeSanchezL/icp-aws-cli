package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/spf13/cobra"
)

func InitCommands(dynamodbClient *dynamodb.Client) *cobra.Command {
	var dynamodbCmd = &cobra.Command{
		Use:   "dynamodb",
		Short: "Commands to interact with Amazon DynamoDB",
		Long:  "Allows listing and managing tables in Amazon DynamoDB.",
	}

	var listTablesCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists DynamoDB tables",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTables(dynamodbClient)
		},
	}

	dynamodbCmd.AddCommand(listTablesCmd)
	return dynamodbCmd
}

func listTables(dynamodbClient *dynamodb.Client) error {
	result, err := dynamodbClient.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		return fmt.Errorf("error listing DynamoDB tables: %w", err)
	}

	for _, tableName := range result.TableNames {
		fmt.Println(tableName)
	}
	return nil
}
