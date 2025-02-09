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
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteObject(s3Client, args[0], args[1])
		},
	}

	deleteBucketCmd := &cobra.Command{
		Use:   "deleteBucket",
		Short: "Deletes an S3 bucket",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteBucket(s3Client, args[0])
		},
	}

	s3Command.AddCommand(deleteObjectCmd)
	s3Command.AddCommand(deleteBucketCmd)
}

func deleteBucket(s3Client *s3.Client, bucketName string) error {
	_, err := s3Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		return fmt.Errorf("failed to delete bucket %s: %v", bucketName, err)
	}

	fmt.Printf("Bucket %s deleted successfully!\n", bucketName)
	return nil
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
