package rds_util

import (
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

type Instance struct{}

func (i Instance) Describe(svc *rds.RDS, name string) ([]string, error) {
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(name),
	}

	result, err := svc.DescribeDBInstances(input)
	if err != nil {
		return nil, err
	}

	instanceInfo := result.DBInstances[0]
	instanceParam := result.DBInstances[0].DBParameterGroups[0]

	values := []*string{
		instanceInfo.DBInstanceIdentifier,
		instanceInfo.EngineVersion,
		instanceInfo.DBInstanceStatus,
		instanceParam.ParameterApplyStatus,
		instanceInfo.PercentProgress,
	}

	var resultValues []string
	for _, value := range values {
		if value != nil {
			resultValues = append(resultValues, *value)
		} else {
			resultValues = append(resultValues, "")
		}
	}

	return resultValues, nil
}

func (i Instance) GetHeaders() []string {
	return []string{"Time", "Duration", "Instance Name", "Version", "Status", "Param Status"}
}

func (i Instance) SaveListToFile(svc *rds.RDS, fileName string) error {
	input := &rds.DescribeDBInstancesInput{}
	result, err := svc.DescribeDBInstances(input)
	if err != nil {
		return err
	}

	var names []string
	for _, instance := range result.DBInstances {
		names = append(names, *instance.DBInstanceIdentifier)
	}

	content := strings.Join(names, "\n")
	return ioutil.WriteFile(fileName, []byte(content), 0644)
}
