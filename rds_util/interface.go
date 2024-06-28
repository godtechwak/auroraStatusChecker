package rds_util

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type Describable interface {
	Describe(svc *rds.RDS, name string) ([]string, error)
	GetHeaders() []string
	SaveListToFile(svc *rds.RDS, fileName string) error
}

func NewRDS(sess *session.Session) *rds.RDS {
	return rds.New(sess)
}
