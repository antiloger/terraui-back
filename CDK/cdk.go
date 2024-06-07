package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfrastructureStackProps struct {
	awscdk.StackProps
}

func UrlShopStack(scope constructs.Construct, id string, props *InfrastructureStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// S3 bucket
	bucket := awss3.NewBucket(stack, jsii.String("urlshop_store"), &awss3.BucketProps{
		RemovalPolicy: awscdk.RemovalPolicy_RETAIN,
	})

	// Database table - [Assumed] : 2 tables with items partiotoned based in users and user_details table to store user data
	//add variable to both table in case of need to apply index
	awsdynamodb.NewTable(stack, jsii.String("items_store"), &awsdynamodb.TableProps{
		TableName: jsii.String("items_store"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("user_id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("item_id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: awscdk.RemovalPolicy_RETAIN,
	})

	//Assumed [To Nawodya] : no need to store tenant data as we have limited set of tenants and they are known (like 3-4)
	awsdynamodb.NewTable(stack, jsii.String("user_details"), &awsdynamodb.TableProps{
		TableName: jsii.String("user_details"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("tenant_id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("user_id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: awscdk.RemovalPolicy_RETAIN,
	})

	// adding golobal sec Indx -- Temp removed this as already implemented the optimized structure for commom queries
	// itemTable.AddLocalSecondaryIndex(&awsdynamodb.LocalSecondaryIndexProps{
	// 	IndexName: jsii.String("CreationDateIndex"),
	// 	SortKey: &awsdynamodb.Attribute{
	// 		Name: jsii.String("price"),
	// 		Type: awsdynamodb.AttributeType_STRING,
	// 	},
	// 	ProjectionType: awsdynamodb.ProjectionType_KEYS_ONLY,
	// })

	// EC2 instance
	vpc := awsec2.NewVpc(stack, jsii.String("VPC"), nil)

	role := awsiam.NewRole(stack, jsii.String("InstanceRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonS3FullAccess")),
		},
	})

	//Used Ubuntu for the VPS
	ubuntuAmi := awsec2.MachineImage_Lookup(&awsec2.LookupMachineImageProps{
		Name:   jsii.String("ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"),
		Owners: &[]*string{jsii.String("099720109477")}, // Ubuntu's owner ID
	})

	instance := awsec2.NewInstance(stack, jsii.String("Instance"), &awsec2.InstanceProps{
		InstanceType: awsec2.InstanceType_Of(awsec2.InstanceClass_T2, awsec2.InstanceSize_MICRO),
		MachineImage: ubuntuAmi,
		Vpc:          vpc,
		Role:         role,
	})

	// Commansd to be executed on the instance and a user data script to get the executable whenever we reboot the instance
	userData := "sudo apt-get update && sudo apt-get install -y nginx awscli\n" +
		"aws s3 cp s3://" + *bucket.BucketName() + "/your-app-binary /usr/local/bin/app-binary\n" +
		"chmod +x /usr/local/bin/app-binary\n" +
		"echo 'server { listen 80; location / { proxy_pass http://localhost:8080; } }' > /etc/nginx/sites-available/default\n" +
		"systemctl restart nginx\n" +
		"echo '@reboot aws s3 cp s3://" + *bucket.BucketName() + "/your-app-binary /usr/local/bin/app-binary' | crontab -\n" +
		"echo '@reboot chmod +x /usr/local/bin/app-binary' | crontab -\n" +
		"echo '@reboot systemctl restart app-binary.service' | crontab -\n"

	instance.AddUserData(jsii.String(userData))

	return stack
}

func main() {
	app := awscdk.NewApp(nil)
	UrlShopStack(app, "UrlShopStack", &InfrastructureStackProps{
		StackProps: awscdk.StackProps{
			Env: &awscdk.Environment{
				Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
				Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
			},
		},
	})
	app.Synth(nil)
}
