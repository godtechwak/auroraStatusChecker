package rds_util

import (
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

type Cluster struct{}

func (c Cluster) Describe(svc *rds.RDS, name string) ([]string, error) {
	input := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(name),
	}

	result, err := svc.DescribeDBClusters(input)
	if err != nil {
		return nil, err
	}

	clusterInfo := result.DBClusters[0]
	clusterParam := result.DBClusters[0].DBClusterMembers[0]

	values := []*string{
		clusterInfo.DBClusterIdentifier,
		clusterInfo.EngineVersion,
		clusterInfo.Status,
		clusterParam.DBClusterParameterGroupStatus,
		clusterInfo.DBClusterInstanceClass,
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

func (c Cluster) GetHeaders() []string {
	return []string{"Time", "Duration", "Cluster Name", "Version", "Status", "Param Status"}
}

func (c Cluster) SaveListToFile(svc *rds.RDS, fileName string) error {
	input := &rds.DescribeDBClustersInput{}
	result, err := svc.DescribeDBClusters(input)
	if err != nil {
		return err
	}

	var names []string
	for _, cluster := range result.DBClusters {
		names = append(names, *cluster.DBClusterIdentifier)
	}

	content := strings.Join(names, "\n")
	return ioutil.WriteFile(fileName, []byte(content), 0644)
}
