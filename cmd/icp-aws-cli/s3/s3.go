package s3

import (
	"icp-aws-cli/cmd/icp-aws-cli/s3/commands"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

func InitCommands(s3Client *s3.Client) *cobra.Command {
	s3Cmd := &cobra.Command{
		Use:   "s3",
		Short: "Commands to interact with Amazon S3",
		Long:  "Allows listing, deleting, and managing objects in Amazon S3.",
	}

	// Initialize subcommands
	commands.InitCopyCommands(s3Client, s3Cmd)
	commands.InitDeleteCommands(s3Client, s3Cmd)
	commands.InitListCommands(s3Client, s3Cmd)
	commands.InitCreateCommands(s3Client, s3Cmd)

	return s3Cmd
}
