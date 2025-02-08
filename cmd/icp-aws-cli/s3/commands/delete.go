package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

func InitDeleteCommands(s3Client *s3.Client, s3Command *cobra.Command) {
	deleteObjectCmd := &cobra.Command{
		Use:   "deleteObject",
		Short: "Deletes a specific object from an S3 bucket",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteObject(s3Client, args[0], args[1])
		},
	}

	s3Command.AddCommand(deleteObjectCmd)
}

func deleteObject(s3Client *s3.Client, bucketName string, objectKey string) error {
	_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	if err != nil {
		return fmt.Errorf("error deleting object %s: %w", objectKey, err)
	}

	fmt.Printf("Object %s deleted successfully from bucket %s\n", objectKey, bucketName)
	return nil
}
