package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

func InitCommands(s3Client *s3.Client) *cobra.Command {
	var s3Cmd = &cobra.Command{
		Use:   "s3",
		Short: "Commands to interact with Amazon S3",
		Long:  "Allows listing and managing buckets in Amazon S3.",
	}

	var listBucketsCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists S3 buckets",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listBuckets(s3Client)
		},
	}

	s3Cmd.AddCommand(listBucketsCmd)
	return s3Cmd
}

func listBuckets(s3Client *s3.Client) error {
	result, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("error listing buckets: %w", err)
	}

	for _, bucket := range result.Buckets {
		fmt.Println(*bucket.Name)
	}
	return nil
}
