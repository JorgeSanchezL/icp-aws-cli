package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type ActionFunc func(*ec2.Client, context.Context, interface{}) (interface{}, error)
type InputBuilderFunc func([]string) interface{}

func manageInstancesWithFilters(ec2Client *ec2.Client, filters []types.Filter, buildInput InputBuilderFunc, actionFunc ActionFunc) error {
	result, err := ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: filters,
	})
	if err != nil {
		return fmt.Errorf("error describing instances: %w", err)
	}

	instanceIDs := []string{}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceIDs = append(instanceIDs, *instance.InstanceId)
		}
	}

	if len(instanceIDs) == 0 {
		return fmt.Errorf("no instances found with the specified filters")
	}

	input := buildInput(instanceIDs)
	_, err = actionFunc(ec2Client, context.TODO(), input)
	if err != nil {
		return fmt.Errorf("error managing instances: %w", err)
	}

	fmt.Printf("Instances %v managed successfully\n", instanceIDs)
	return nil
}
