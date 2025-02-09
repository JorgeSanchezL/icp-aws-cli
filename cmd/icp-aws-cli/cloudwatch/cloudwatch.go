package cloudwatch

import (
	"icp-aws-cli/cmd/icp-aws-cli/cloudwatch/commands/alarms"
	"icp-aws-cli/cmd/icp-aws-cli/cloudwatch/commands/logs/events"
	"icp-aws-cli/cmd/icp-aws-cli/cloudwatch/commands/logs/loggroups"
	"icp-aws-cli/cmd/icp-aws-cli/cloudwatch/commands/logs/streams"
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
	alarms.InitDeleteAlarmCommand(cwClient, cloudWatchCmd)
	loggroups.InitCreateLogGroupCommand(cwLogsClient, cloudWatchCmd)
	loggroups.InitListLogGroupsCommand(cwLogsClient, cloudWatchCmd)
	loggroups.InitDeleteLogGroupCommand(cwLogsClient, cloudWatchCmd)
	events.InitGetLogEventsCommand(cwLogsClient, cloudWatchCmd)
	streams.InitListLogStreamsCommand(cwLogsClient, cloudWatchCmd)
	metrics.InitCreateMetricCommand(cwClient, cloudWatchCmd)
	metrics.InitListMetricsCommand(cwClient, cloudWatchCmd)
	metrics.InitDeleteMetricCommand(cwClient, cloudWatchCmd)

	return cloudWatchCmd
}
