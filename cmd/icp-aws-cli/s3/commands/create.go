package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

func InitCreateCommands(s3Client *s3.Client, s3Command *cobra.Command) {
	createBucketCmd := &cobra.Command{
		Use:   "createBucket",
		Short: "Creates a new S3 bucket",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createBucket(s3Client, args[0])
		},
	}

	s3Command.AddCommand(createBucketCmd)
}

func createBucket(s3Client *s3.Client, bucketName string) error {
	_, err := s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %v", bucketName, err)
	}

	fmt.Printf("Bucket %s created successfully!\n", bucketName)
	return nil
}
