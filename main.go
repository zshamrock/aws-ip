package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	s "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/chyeh/pubip"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"strings"
)

const (
	groupNameFlagName = "group-name"
	portFlagName      = "port"
	locationFlagName  = "location"

	groupNameValueSeparator = ","
)

const (
	appName = "aws-ip"
	version = "1.0.0"
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Usage = `Sync/set AWS security group entry (by description) with current user's local public IP address`
	app.Version = version
	app.Author = "(c) Aliaksandr Kazlou"
	app.Metadata = map[string]interface{}{"GitHub": "https://github.com/zshamrock/aws-ip"}
	app.UsageText = fmt.Sprintf(`%s		 
        --group-name    <comma separated affected EC2 security groups> 
        --port          <port>	
        --location      <free text/code of the current user's location, like home, office, coworking, etc.>`,
		appName)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  fmt.Sprintf("%s", groupNameFlagName),
			Usage: "comma separated affected EC2 security groups",
		},
		cli.Int64Flag{
			Name:  fmt.Sprintf("%s", portFlagName),
			Usage: "port number",
		},
		cli.StringFlag{
			Name:  fmt.Sprintf("%s", locationFlagName),
			Usage: "free text/code of the current user's location, like home, office, coworking, etc.",
		},
	}
	app.Action = action

	err := app.Run(os.Args)
	if err != nil {
		log.Panicf("error encountered while running the app %v", err)
	}
}

func action(c *cli.Context) error {
	if len(os.Args) == 1 {
		cli.ShowAppHelpAndExit(c, 0)
	}
	groups := strings.Split(mustStringFlag(c, groupNameFlagName), groupNameValueSeparator)
	port := mustInt64Flag(c, portFlagName)
	location := mustStringFlag(c, locationFlagName)
	session, err := s.NewSessionWithOptions(s.Options{
		SharedConfigState: s.SharedConfigEnable,
	})
	if err != nil {
		log.Panic("unable to establish AWS session connection: ", err)
	}
	username, err := getUserName(session)
	if err != nil {
		log.Panic("unable to determine implicit username by used AWS profile")
		return err
	}
	ipAddress := findIPAddress()
	for _, group := range groups {
		err := syncSecurityGroupInboundRule(session, group, port, ipAddress, username, location)
		if err != nil {
			return err
		}
	}
	return nil
}

func findIPAddress() string {
	ip, err := pubip.Get()
	if err != nil {
		log.Panic("couldn't get host public IP address: ", err)
	}
	return ip.To4().String()
}

func getUserName(session *s.Session) (string, error) {
	svc := iam.New(session)
	result, err := svc.GetUser(&iam.GetUserInput{})
	if err != nil {
		return "", err
	}
	return aws.StringValue(result.User.UserName), nil
}

func syncSecurityGroupInboundRule(
	session *s.Session, groupName string, port int64, ipAddress string, username string, location string) error {
	svc := ec2.New(session)
	groups, err := findSecurityGroups(svc, groupName)
	if err != nil {
		return err
	}
	descriptionId := buildDescriptionId(username, location)
	for _, group := range groups {
		revoked := false
		for _, inbound := range group.IpPermissions {
			if revoked {
				break
			}
			for _, entry := range inbound.IpRanges {
				if aws.StringValue(entry.Description) == descriptionId {
					err := revokeSecurityGroupIngress(svc, groupName, port, aws.StringValue(entry.CidrIp), username, location)
					if err != nil {
						return err
					}
					revoked = true
					break
				}
			}
		}
		err = authorizeSecurityGroupIngress(svc, groupName, port, fmt.Sprintf("%s/32", ipAddress), username, location)
		if err != nil {
			return err
		}
	}
	return nil
}

func findSecurityGroups(svc *ec2.EC2, groupName string) ([]*ec2.SecurityGroup, error) {
	fmt.Print(":: find security groups...")
	input := &ec2.DescribeSecurityGroupsInput{
		GroupNames: []*string{
			aws.String(groupName),
		},
	}

	result, err := svc.DescribeSecurityGroups(input)
	if err != nil {
		return nil, err
	}
	fmt.Println("done")
	return result.SecurityGroups, nil
}

func revokeSecurityGroupIngress(svc *ec2.EC2, groupName string, port int64, ipAddress string, username string, location string) error {
	fmt.Print(":: revoking security group ingress...")
	input := &ec2.RevokeSecurityGroupIngressInput{
		GroupName: aws.String(groupName),
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort:   aws.Int64(port),
				ToPort:     aws.Int64(port),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String(ipAddress),
						Description: aws.String(buildDescriptionId(username, location)),
					},
				},
			},
		},
	}
	_, err := svc.RevokeSecurityGroupIngress(input)
	if err != nil {
		return err
	}
	fmt.Println("done")
	return nil
}

func authorizeSecurityGroupIngress(
	svc *ec2.EC2, groupName string, port int64, ipAddress string, username string, location string) error {
	fmt.Print(":: authorizing security group ingress...")
	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupName: aws.String(groupName),
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort:   aws.Int64(port),
				ToPort:     aws.Int64(port),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String(ipAddress),
						Description: aws.String(buildDescriptionId(username, location)),
					},
				},
			},
		},
	}
	_, err := svc.AuthorizeSecurityGroupIngress(input)
	if err != nil {
		return err
	}
	fmt.Println("done")
	return nil
}

func buildDescriptionId(username string, location string) string {
	return fmt.Sprintf("%s-%s", username, location)
}

func mustStringFlag(c *cli.Context, name string) string {
	value := c.String(name)
	if value == "" {
		log.Panic(fmt.Sprintf("%s is required", name))
	}
	return value
}

func mustInt64Flag(c *cli.Context, name string) int64 {
	value := c.Int64(name)
	if value == 0 {
		log.Panic(fmt.Sprintf("%s is required", name))
	}
	return value
}
