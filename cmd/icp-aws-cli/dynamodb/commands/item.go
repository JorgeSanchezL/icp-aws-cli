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

func InitItemCommands(dynamodbClient *dynamodb.Client, dynamodbCmd *cobra.Command) {
	putItemCmd := &cobra.Command{
		Use:   "putItem",
		Short: "Puts an item into a DynamoDB table",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return putItem(dynamodbClient, args[0], args[1])
		},
	}

	getItemCmd := &cobra.Command{
		Use:   "getItem",
		Short: "Retrieves an item from a DynamoDB table",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getItem(dynamodbClient, args[0], args[1])
		},
	}

	dynamodbCmd.AddCommand(putItemCmd)
	dynamodbCmd.AddCommand(getItemCmd)
}

// putItems adds a new item, given as JSON, to the given table name
func putItem(client *dynamodb.Client, tableName string, itemJSON string) error {
	var item map[string]interface{}
	if err := json.Unmarshal([]byte(itemJSON), &item); err != nil {
		return fmt.Errorf("error parsing item JSON: %w", err)
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("error marshaling item: %w", err)
	}
	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("error putting item: %w", err)
	}

	fmt.Printf("Item put successfully into table %s\n", tableName)
	return nil
}

// getItem retrieves an item from the provided table name
func getItem(client *dynamodb.Client, tableName string, keyJSON string) error {
	var key map[string]interface{}
	if err := json.Unmarshal([]byte(keyJSON), &key); err != nil {
		return fmt.Errorf("error parsing key JSON: %w", err)
	}

	av, err := attributevalue.MarshalMap(key)
	if err != nil {
		return fmt.Errorf("error marshaling key: %w", err)
	}

	result, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       av,
	})
	if err != nil {
		return fmt.Errorf("error getting item: %w", err)
	}

	if len(result.Item) == 0 {
		fmt.Println("No item found")
		return nil
	}

	var item map[string]interface{}
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return fmt.Errorf("error unmarshaling item: %w", err)
	}

	fmt.Println(item)
	return nil
}
