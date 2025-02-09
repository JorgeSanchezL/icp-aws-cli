package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

func InitCopyCommands(s3Client *s3.Client, s3Command *cobra.Command) {
	copyObjectCmd := &cobra.Command{
		Use:   "copyObject",
		Short: "Copies an object from one bucket to another",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			return copyObject(s3Client, args[0], args[1], args[2], args[3])
		},
	}

	s3Command.AddCommand(copyObjectCmd)
}

func copyObject(s3Client *s3.Client, srcBucket string, srcKey string, destBucket string, destKey string) error {
	copySource := fmt.Sprintf("%s/%s", srcBucket, srcKey)
	_, err := s3Client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     &destBucket,
		CopySource: &copySource,
		Key:        &destKey,
	})
	if err != nil {
		return fmt.Errorf("error copying object %s from bucket %s to %s: %w", srcKey, srcBucket, destBucket, err)
	}

	fmt.Printf("Object %s copied successfully from bucket %s to bucket %s as %s\n", srcKey, srcBucket, destBucket, destKey)
	return nil
}
