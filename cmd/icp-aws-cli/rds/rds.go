package rds

import (
	"icp-aws-cli/cmd/icp-aws-cli/rds/commands"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/spf13/cobra"
)

func InitCommands(rdsClient *rds.Client) *cobra.Command {
	rdsCmd := &cobra.Command{
		Use:   "rds",
		Short: "Commands to interact with Amazon RDS",
		Long:  "Allows listing and managing RDS instances in Amazon RDS.",
	}

	// Initialize subcommands
	commands.InitListCommands(rdsClient, rdsCmd)
	commands.InitCreateCommands(rdsClient, rdsCmd)
	commands.InitDeleteCommands(rdsClient, rdsCmd)
	commands.InitStartStopCommands(rdsClient, rdsCmd)

	return rdsCmd
}
