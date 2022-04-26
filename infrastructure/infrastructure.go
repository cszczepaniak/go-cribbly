package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	BuildHash         = `BUILD_HASH`
	ApplicationBucket = `CRIBBLY_APP_BUCKET`
)

func lookup(k string) string {
	v, ok := os.LookupEnv(k)
	if !ok {
		panic(fmt.Sprintf(`environment variable %q not found`, k))
	}
	return v
}

type CDKEnvironment struct {
	buildHash             string
	applicationBucketName string
}

func getCDKEnv() (CDKEnvironment, error) {
	return CDKEnvironment{
		buildHash:             lookup(BuildHash),
		applicationBucketName: lookup(ApplicationBucket),
	}, nil
}

type ApplicationStackProps struct {
	baseProps           awscdk.StackProps
	infrastructureStack InfrastructureStack
	env                 CDKEnvironment
}

// NewApplicationStack configures the stack for the application.
func NewApplicationStack(scope constructs.Construct, id string, props *ApplicationStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.baseProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	stack.AddDependency(props.infrastructureStack.baseStack, jsii.String(`application stack depends on the infrastructure stack`))

	// Resource: S3 bucket for storing application data
	dataBucket := awss3.NewBucket(stack, inheritID(id, `data-bucket`), &awss3.BucketProps{
		BucketName:       jsii.String(`cribbly-data-bucket`),
		PublicReadAccess: jsii.Bool(false),
		Versioned:        jsii.Bool(false),
	})

	// Resource: the lambda handler
	lambda := awslambda.NewFunction(stack, inheritID(id, `lambda`), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Handler: jsii.String(`cribbly-backend`),
		Code:    awslambda.AssetCode_FromBucket(props.infrastructureStack.appBucket, jsii.String(props.env.buildHash+`.zip`), nil),
		Environment: &map[string]*string{
			`GIN_MODE`:            jsii.String(`release`),
			`CRIBBLY_DATA_BUCKET`: dataBucket.BucketName(),
		},
	})

	dataBucket.GrantReadWrite(lambda, nil)

	// Resource: the API gateway that dispatches incoming requests to the lambda
	awsapigateway.NewLambdaRestApi(stack, inheritID(id, `apigateway`), &awsapigateway.LambdaRestApiProps{
		RestApiName: jsii.String(`cribbly-api`),
		Handler:     lambda,
		DeployOptions: &awsapigateway.StageOptions{
			StageName: jsii.String(`dev`),
		},
	})

	return stack
}

type InfrastructureStackProps struct {
	baseProps awscdk.StackProps
	env       CDKEnvironment
}

type InfrastructureStack struct {
	baseStack awscdk.Stack
	appBucket awss3.Bucket
}

// NewInfrastructureStack configures the stack that is needed to host infrastructure artifacts and resources.
func NewInfrastructureStack(scope constructs.Construct, id string, props *InfrastructureStackProps) InfrastructureStack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.baseProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	lambdaAppBucket := awss3.NewBucket(stack, inheritID(id, `app-bucket`), &awss3.BucketProps{
		BucketName: jsii.String(props.env.applicationBucketName),
		Versioned:  jsii.Bool(false),
	})

	return InfrastructureStack{
		baseStack: stack,
		appBucket: lambdaAppBucket,
	}
}

func main() {
	envVars, err := getCDKEnv()
	if err != nil {
		panic(err)
	}

	app := awscdk.NewApp(nil)

	infraStack := NewInfrastructureStack(app, `CribblyInfrastructureStack`, &InfrastructureStackProps{
		baseProps: awscdk.StackProps{
			Env: env(),
		},
		env: envVars,
	})

	NewApplicationStack(app, `CribblyApplicationStack`, &ApplicationStackProps{
		baseProps: awscdk.StackProps{
			Env: env(),
		},
		infrastructureStack: infraStack,
		env:                 envVars,
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(`977779885689`),
		Region:  jsii.String(`us-east-2`),
	}
}

func inheritID(base string, id string) *string {
	return jsii.String(base + `-` + id)
}
