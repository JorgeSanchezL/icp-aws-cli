package dynamodb

import (
	"icp-aws-cli/cmd/icp-aws-cli/dynamodb/commands"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/spf13/cobra"
)

func InitCommands(dynamodbClient *dynamodb.Client) *cobra.Command {
	dynamodbCmd := &cobra.Command{
		Use:   "dynamodb",
		Short: "Commands to interact with Amazon DynamoDB",
		Long:  "Allows listing and managing tables in Amazon DynamoDB.",
	}

	// Initialize subcommands
	commands.InitCreateCommands(dynamodbClient, dynamodbCmd)
	commands.InitDeleteCommands(dynamodbClient, dynamodbCmd)
	commands.InitDescribeCommands(dynamodbClient, dynamodbCmd)
	commands.InitItemCommands(dynamodbClient, dynamodbCmd)
	commands.InitListCommands(dynamodbClient, dynamodbCmd)
	commands.InitQueryCommands(dynamodbClient, dynamodbCmd)

	return dynamodbCmd
}
