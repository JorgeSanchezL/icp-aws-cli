package ec2

import (
	"icp-aws-cli/cmd/icp-aws-cli/ec2/commands"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
)

func InitCommands(ec2Client *ec2.Client) *cobra.Command {
	var ec2Cmd = &cobra.Command{
		Use:   "ec2",
		Short: "Commands to interact with Amazon EC2",
		Long:  "Allows listing and managing EC2 instances in Amazon EC2.",
	}

	// Initialize subcommands
	commands.InitListCommands(ec2Client, ec2Cmd)
	commands.InitStartCommands(ec2Client, ec2Cmd)
	commands.InitStopCommands(ec2Client, ec2Cmd)
	commands.InitRebootCommands(ec2Client, ec2Cmd)
	commands.InitTerminateCommands(ec2Client, ec2Cmd)
	commands.InitCreateCommands(ec2Client, ec2Cmd)

	return ec2Cmd
}
