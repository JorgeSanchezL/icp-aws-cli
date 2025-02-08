package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

func InitListCommands(s3Client *s3.Client, s3Command *cobra.Command) {
	listBucketsCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists S3 buckets",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listBuckets(s3Client)
		},
	}

	listObjectsCmd := &cobra.Command{
		Use:   "listObjects",
		Short: "Lists all objects of an S3 bucket",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listObjects(s3Client, args[0])
		},
	}

	listObjectsByExtensionCmd := &cobra.Command{
		Use:   "listObjectsByExtension",
		Short: "Lists all objects in a bucket with a specific file extension",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listObjectsByExtension(s3Client, args[0], args[1])
		},
	}

	s3Command.AddCommand(listBucketsCmd)
	s3Command.AddCommand(listObjectsCmd)
	s3Command.AddCommand(listObjectsByExtensionCmd)
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

func listObjects(s3Client *s3.Client, bucketName string) error {
	result, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{Bucket: &bucketName})
	if err != nil {
		return fmt.Errorf("error listing objects: %w", err)
	}

	for _, object := range result.Contents {
		fmt.Println(*object.Key)
	}
	return nil
}

func listObjectsByExtension(s3Client *s3.Client, bucketName string, extension string) error {
	result, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{Bucket: &bucketName})
	if err != nil {
		return fmt.Errorf("error listing objects: %w", err)
	}

	for _, object := range result.Contents {
		if strings.HasSuffix(*object.Key, "."+extension) {
			fmt.Println(*object.Key)
		}
	}
	return nil
}
