package s3

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

func InitCommands(s3Client *s3.Client) *cobra.Command {
	s3Cmd := &cobra.Command{
		Use:   "s3",
		Short: "Commands to interact with Amazon S3",
		Long:  "Allows listing, deleting, and managing objects in Amazon S3.",
	}

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

	deleteObjectCmd := &cobra.Command{
		Use:   "deleteObject",
		Short: "Deletes a specific object from an S3 bucket",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteObject(s3Client, args[0], args[1])
		},
	}

	listObjectsByExtensionCmd := &cobra.Command{
		Use:   "listObjectsByExtension",
		Short: "Lists all objects in a bucket with a specific file extension",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listObjectsByExtension(s3Client, args[0], args[1])
		},
	}

	copyObjectCmd := &cobra.Command{
		Use:   "copyObject",
		Short: "Copies an object from one bucket to another",
		RunE: func(cmd *cobra.Command, args []string) error {
			return copyObject(s3Client, args[0], args[1], args[2], args[3])
		},
	}

	s3Cmd.AddCommand(listBucketsCmd)
	s3Cmd.AddCommand(listObjectsCmd)
	s3Cmd.AddCommand(deleteObjectCmd)
	s3Cmd.AddCommand(listObjectsByExtensionCmd)
	s3Cmd.AddCommand(copyObjectCmd)
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
