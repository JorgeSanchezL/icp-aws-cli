package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func getInstanceIDsByPattern(ec2Client *ec2.Client, pattern string) ([]string, error) {
	result, err := ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("error describing instances: %w", err)
	}

	var instanceIDs []string
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if strings.Contains(*instance.InstanceId, pattern) {
				instanceIDs = append(instanceIDs, *instance.InstanceId)
			}
		}
	}
	return instanceIDs, nil
}

func getInstanceIDsByTag(ec2Client *ec2.Client, key, value string) ([]string, error) {
	result, err := ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String(fmt.Sprintf("tag:%s", key)),
				Values: []string{value},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error describing instances by tag: %w", err)
	}

	var instanceIDs []string
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceIDs = append(instanceIDs, *instance.InstanceId)
		}
	}
	return instanceIDs, nil
}
