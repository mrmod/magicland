package magicland

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var ec2Client *ec2.EC2

func init() {
	// Leverages AWS_ACCESS_KEY_ID || AWS_ACCESS_KEY
	// And AWS_SECRET_ACCESS_KEY || AWS_SECRET_KEY
	config := aws.NewConfig().WithCredentials(credentials.NewEnvCredentials())
	config = config.WithRegion(envOr("AWS_REGION", "us-west-2"))

	ec2Client = ec2.New(session.New(config))
}

// listInstances Nothing yet
func listInstances() {
	instances, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		fmt.Println("Failed to describe instances", err)
		return
	}

	fmt.Println(instances)
}
