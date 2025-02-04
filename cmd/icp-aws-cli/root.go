package commands

import (
	"icp-aws-cli/cmd/icp-aws-cli/autoscaling"
	"icp-aws-cli/cmd/icp-aws-cli/cloudwatch"
	"icp-aws-cli/cmd/icp-aws-cli/dynamodb"
	"icp-aws-cli/cmd/icp-aws-cli/ec2"
	"icp-aws-cli/cmd/icp-aws-cli/rds"
	"icp-aws-cli/cmd/icp-aws-cli/s3"
	"icp-aws-cli/pkg/awsclient"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "icp-aws-cli",
	Short: "CLI to interact with AWS",
	Long:  "A CLI in Go to manage AWS resources like EC2, S3, DynamoDB, AutoScaling, RDS, and CloudWatch.",
}

func InitCommands(clients *awsclient.AWSClientCollection) {
	RootCmd.AddCommand(s3.InitCommands(clients.S3))
	RootCmd.AddCommand(ec2.InitCommands(clients.EC2))
	RootCmd.AddCommand(dynamodb.InitCommands(clients.DynamoDB))
	RootCmd.AddCommand(rds.InitCommands(clients.RDS))
	RootCmd.AddCommand(cloudwatch.InitCommands(clients.CloudWatch))
	RootCmd.AddCommand(autoscaling.InitCommands(clients.AutoScaling))
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}
