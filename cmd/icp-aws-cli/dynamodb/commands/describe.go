package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/spf13/cobra"
)

func InitDescribeCommands(dynamodbClient *dynamodb.Client, dynamodbCmd *cobra.Command) {
	describeTableCmd := &cobra.Command{
		Use:   "describe",
		Short: "Describes a DynamoDB table",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return describeTable(dynamodbClient, args[0])
		},
	}

	dynamodbCmd.AddCommand(describeTableCmd)
}

// describeTable describes a DynamoDB table
func describeTable(client *dynamodb.Client, tableName string) error {
	result, err := client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return fmt.Errorf("error describing table %s: %w", tableName, err)
	}

	fmt.Printf("Table Name: %s\n", *result.Table.TableName)
	fmt.Printf("Status: %s\n", result.Table.TableStatus)
	fmt.Printf("Item Count: %d\n", result.Table.ItemCount)
	return nil
}
