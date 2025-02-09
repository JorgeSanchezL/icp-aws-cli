package cloudwatch

import (
	"icp-aws-cli/cmd/icp-aws-cli/cloudwatch/commands/alarms"
	"icp-aws-cli/cmd/icp-aws-cli/cloudwatch/commands/logs"
	"icp-aws-cli/cmd/icp-aws-cli/cloudwatch/commands/metrics"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

func InitCommands(cwClient *cloudwatch.Client, cwLogsClient *cloudwatchlogs.Client) *cobra.Command {
	var cloudWatchCmd = &cobra.Command{
		Use:   "cloudwatch",
		Short: "Commands for interacting with AWS CloudWatch",
	}

	// Initialize subcommands
	alarms.InitCreateAlarmCommand(cwClient, cloudWatchCmd)
	alarms.InitListAlarmsCommand(cwClient, cloudWatchCmd)
	alarms.InitDescribeAlarmCommand(cwClient, cloudWatchCmd)
	alarms.InitDeleteAlarmCommand(cwClient, cloudWatchCmd)
	logs.InitListLogGroupsCommand(cwLogsClient, cloudWatchCmd)
	logs.InitDescribeLogGroupCommand(cwLogsClient, cloudWatchCmd)
	metrics.InitListMetricsCommand(cwClient, cloudWatchCmd)

	return cloudWatchCmd
}
