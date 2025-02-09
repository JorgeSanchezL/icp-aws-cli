package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/spf13/cobra"
)

func InitCreateCommands(dynamodbClient *dynamodb.Client, dynamodbCmd *cobra.Command) {
	createTableCmd := &cobra.Command{
		Use:   "createTable",
		Short: "Creates a new DynamoDB table",
		Args:  cobra.RangeArgs(3, 5),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 && len(args) != 5 {
				return fmt.Errorf("createTable requires 3 or 5 arguments")
			}
			skName, skType := "", ""
			if len(args) == 5 {
				skName = args[3]
				skType = args[4]
			}
			return createTable(dynamodbClient, args[0], args[1], args[2], skName, skType)
		},
	}

	dynamodbCmd.AddCommand(createTableCmd)
}

// createTable creates a new DynamoDB table
func createTable(client *dynamodb.Client, tableName, pkName, pkType, skName, skType string) error {
	attrs := []types.AttributeDefinition{{
		AttributeName: aws.String(pkName),
		AttributeType: types.ScalarAttributeType(pkType),
	}}

	keySchema := []types.KeySchemaElement{{
		AttributeName: aws.String(pkName),
		KeyType:       types.KeyTypeHash,
	}}

	if skName != "" && skType != "" {
		attrs = append(attrs, types.AttributeDefinition{
			AttributeName: aws.String(skName),
			AttributeType: types.ScalarAttributeType(skType),
		})
		keySchema = append(keySchema, types.KeySchemaElement{
			AttributeName: aws.String(skName),
			KeyType:       types.KeyTypeRange,
		})
	}
	_, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName:            aws.String(tableName),
		AttributeDefinitions: attrs,
		KeySchema:            keySchema,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})
	if err != nil {
		return fmt.Errorf("error creating table %s: %w", tableName, err)
	}

	fmt.Printf("Table %s created successfully\n", tableName)
	return nil
}
