package awscfg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func Connect() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String(`us-east-2`),
		Credentials: credentials.NewEnvCredentials(),
	})
}
