package awsclient

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSClientCollection struct {
	S3          *s3.Client
	EC2         *ec2.Client
	DynamoDB    *dynamodb.Client
	AutoScaling *autoscaling.Client
	RDS         *rds.Client
	CloudWatch  *cloudwatch.Client
}

func NewAWSClientCollection() (*AWSClientCollection, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error loading AWS config: %w", err)
	}

	return &AWSClientCollection{
		S3:          s3.NewFromConfig(cfg),
		EC2:         ec2.NewFromConfig(cfg),
		DynamoDB:    dynamodb.NewFromConfig(cfg),
		AutoScaling: autoscaling.NewFromConfig(cfg),
		RDS:         rds.NewFromConfig(cfg),
		CloudWatch:  cloudwatch.NewFromConfig(cfg),
	}, nil
}
