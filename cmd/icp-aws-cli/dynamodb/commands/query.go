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

func InitQueryCommands(dynamodbClient *dynamodb.Client, dynamodbCmd *cobra.Command) {
	queryItemsCmd := &cobra.Command{
		Use:   "query",
		Short: "Queries items in a DynamoDB table",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return queryItems(dynamodbClient, args[0], args[1], args[2])
		},
	}

	dynamodbCmd.AddCommand(queryItemsCmd)
}

func queryItems(client *dynamodb.Client, tableName, keyCondition, exprAttrValuesJSON string) error {
	var exprAttrValues map[string]interface{}
	if err := json.Unmarshal([]byte(exprAttrValuesJSON), &exprAttrValues); err != nil {
		return fmt.Errorf("error parsing expression attribute values JSON: %w", err)
	}

	avs, err := attributevalue.MarshalMap(exprAttrValues)
	if err != nil {
		return fmt.Errorf("error marshaling expression attribute values: %w", err)
	}

	result, err := client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String(keyCondition),
		ExpressionAttributeValues: avs,
	})
	if err != nil {
		return fmt.Errorf("error querying items: %w", err)
	}

	for _, item := range result.Items {
		var itemMap map[string]interface{}
		if err := attributevalue.UnmarshalMap(item, &itemMap); err != nil {
			return fmt.Errorf("error unmarshaling item: %w", err)
		}
		fmt.Println(itemMap)
	}
	return nil
}
