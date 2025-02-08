package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/spf13/cobra"
)

func InitDeleteCommands(dynamodbClient *dynamodb.Client, dynamodbCmd *cobra.Command) {
	deleteTableCmd := &cobra.Command{
		Use:   "deleteTable",
		Short: "Deletes a DynamoDB table",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteTable(dynamodbClient, args[0])
		},
	}

	deleteItemCmd := &cobra.Command{
		Use:   "deleteItem",
		Short: "Deletes an item from a DynamoDB table",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteItem(dynamodbClient, args[0], args[1])
		},
	}

	dynamodbCmd.AddCommand(deleteItemCmd)
	dynamodbCmd.AddCommand(deleteTableCmd)
}

// deleteTable deletes a DynamoDB table
func deleteTable(client *dynamodb.Client, tableName string) error {
	_, err := client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return fmt.Errorf("error deleting table %s: %w", tableName, err)
	}

	fmt.Printf("Table %s deleted successfully\n", tableName)
	return nil
}

// deleteItem deletes an item from the given table name
func deleteItem(client *dynamodb.Client, tableName string, keyJSON string) error {
	var key map[string]interface{}
	if err := json.Unmarshal([]byte(keyJSON), &key); err != nil {
		return fmt.Errorf("error parsing key JSON: %w", err)
	}

	av, err := attributevalue.MarshalMap(key)
	if err != nil {
		return fmt.Errorf("error marshaling key: %w", err)
	}

	_, err = client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       av,
	})
	if err != nil {
		return fmt.Errorf("error deleting item: %w", err)
	}

	fmt.Printf("Item deleted successfully from table %s\n", tableName)
	return nil
}
