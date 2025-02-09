package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/spf13/cobra"
)

func InitListCommands(dynamodbClient *dynamodb.Client, dynamodbCmd *cobra.Command) {
	listTablesCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists DynamoDB tables",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTables(dynamodbClient)
		},
	}

	dynamodbCmd.AddCommand(listTablesCmd)
}

// listTables retrieves all the DynamoDB tables the current user has access to
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
