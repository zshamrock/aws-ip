package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
)

func main() {
	username, err := getUserName()
	if err != nil {
		fmt.Println("Unable to determine username implicitly by used AWS profile")
		return
	}
	fmt.Println(username)
	getSecurityGroupInboundRule(username, "home")
}

func getUserName() (string, error) {
	svc := iam.New(session.New())
	result, err := svc.GetUser(&iam.GetUserInput{})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", err
	}

	return aws.StringValue(result.User.UserName), nil
}

func getSecurityGroupInboundRule(username string, location string) (string, error) {
	session, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	svc := ec2.New(session)
	input := &ec2.DescribeSecurityGroupsInput{
		GroupNames: []*string{
			aws.String(""),
		},
	}

	result, err := svc.DescribeSecurityGroups(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", err
	}

	for _, sg := range result.SecurityGroups {
		for _, inbound := range sg.IpPermissions {
			for _, entry := range inbound.IpRanges {
				if aws.StringValue(entry.Description) == fmt.Sprintf("%s-%s", username, location) {
					fmt.Printf("Found %s\n", entry)
				}
			}
		}
	}

	return "", nil
}
