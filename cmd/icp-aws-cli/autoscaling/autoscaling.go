package autoscaling

import (
	"icp-aws-cli/cmd/icp-aws-cli/autoscaling/commands"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/spf13/cobra"
)

func InitCommands(asClient *autoscaling.Client) *cobra.Command {
	var autoscalingCmd = &cobra.Command{
		Use:   "autoscaling",
		Short: "Commands to interact with Amazon AutoScaling",
		Long:  "Allows listing and managing AutoScaling groups in Amazon AutoScaling.",
	}
	// Initialize subcommands
	commands.InitListGroupsCommand(asClient, autoscalingCmd)
	commands.InitDescribeGroupCommand(asClient, autoscalingCmd)
	commands.InitCreateGroupCommand(asClient, autoscalingCmd)
	commands.InitDeleteGroupCommand(asClient, autoscalingCmd)
	commands.InitUpdateGroupCommand(asClient, autoscalingCmd)

	return autoscalingCmd
}
