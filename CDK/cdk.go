package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfrastructureStackProps struct {
	awscdk.StackProps
}

func UrlShopStack(scope constructs.Construct, id string, props *InfrastructureStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// S3 bucket
	// awss3.NewBucket(stack, jsii.String("urlshop_StaticStore"), &awss3.BucketProps{
	// 	RemovalPolicy: awscdk.RemovalPolicy_RETAIN,
	// })

	// // Database table - [Assumed] : 2 tables with items partiotoned based in users and user_details table to store user data
	// //add variable to both table in case of need to apply index
	// awsdynamodb.NewTable(stack, jsii.String("items_store"), &awsdynamodb.TableProps{
	// 	TableName: jsii.String("items_store"),
	// 	PartitionKey: &awsdynamodb.Attribute{
	// 		Name: jsii.String("user_id"),
	// 		Type: awsdynamodb.AttributeType_STRING,
	// 	},
	// 	SortKey: &awsdynamodb.Attribute{
	// 		Name: jsii.String("item_id"),
	// 		Type: awsdynamodb.AttributeType_STRING,
	// 	},
	// 	BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
	// 	RemovalPolicy: awscdk.RemovalPolicy_RETAIN,
	// })

	// //Assumed [To Nawodya] : no need to store tenant data as we have limited set of tenants and they are known (like 3-4)
	// awsdynamodb.NewTable(stack, jsii.String("user_details"), &awsdynamodb.TableProps{
	// 	TableName: jsii.String("user_details"),
	// 	PartitionKey: &awsdynamodb.Attribute{
	// 		Name: jsii.String("tenant_id"),
	// 		Type: awsdynamodb.AttributeType_STRING,
	// 	},
	// 	SortKey: &awsdynamodb.Attribute{
	// 		Name: jsii.String("user_id"),
	// 		Type: awsdynamodb.AttributeType_STRING,
	// 	},
	// 	BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
	// 	RemovalPolicy: awscdk.RemovalPolicy_RETAIN,
	// })

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
	vpc := awsec2.NewVpc(stack, jsii.String("VPC"), &awsec2.VpcProps{
		MaxAzs: jsii.Number(2), // Default is all AZs in region
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:       jsii.String("PublicSubnet"),
				SubnetType: awsec2.SubnetType_PUBLIC,
			},
			// {
			// 	Name:       jsii.String("PrivateSubnet"),
			// 	SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED,
			// },
		},
		NatGateways: jsii.Number(0),
	})

	role := awsiam.NewRole(stack, jsii.String("InstanceRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonS3FullAccess")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonSSMManagedInstanceCore")),
		},
	})

	//Used Ubuntu for the VPS
	ubuntuAmi := awsec2.MachineImage_GenericLinux(&map[string]*string{
		os.Getenv("CDK_DEFAULT_REGION"): jsii.String("ami-003c463c8207b4dfa"),
	}, nil)

	securityGroup := awsec2.NewSecurityGroup(stack, jsii.String("InstanceSG"), &awsec2.SecurityGroupProps{
		Vpc: vpc,
	})
	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(80)), jsii.String("Allow HTTP"), nil)
	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(443)), jsii.String("Allow HTTPS"), nil)

	instance := awsec2.NewInstance(stack, jsii.String("Instance"), &awsec2.InstanceProps{
		InstanceType: awsec2.InstanceType_Of(awsec2.InstanceClass_T2, awsec2.InstanceSize_MICRO),
		MachineImage: ubuntuAmi,
		Vpc:          vpc,
		Role:         role,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PUBLIC,
		},
		SecurityGroup: securityGroup,
	})

	// Commansd to be executed on the instance and a user data script to get the executable whenever we reboot the instance
	userData := `#!/bin/bash

	# Log output to a file and also print to console
	exec > >(tee -a /var/log/user-data.log | logger -t user-data -s 2>/dev/console) 2>&1
	
	# Update the package list
	echo "Updating package list..."
	if sudo apt-get update; then
	  echo "Package list updated successfully."
	else
	  echo "Failed to update package list" >&2
	  exit 1
	fi
	
	# Install necessary packages
	echo "Installing necessary packages..."
	if sudo apt-get install -y nginx bison make build-essential unzip; then
	  echo "Necessary packages installed successfully."
	else
	  echo "Failed to install necessary packages" >&2
	  exit 1
	fi
	
	# Update GLIBC
	echo "Updating GLIBC..."
	if sudo apt-get install -y libc6; then
	  echo "GLIBC updated successfully."
	else
	  echo "Failed to update GLIBC" >&2
	  exit 1
	fi
	
	# Install AWS CLI v2
	echo "Installing AWS CLI v2..."
	if curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"; then
	  if unzip awscliv2.zip; then
		if sudo ./aws/install; then
		  echo "AWS CLI v2 installed successfully."
		else
		  echo "Failed to install AWS CLI v2" >&2
		  exit 1
		fi
		else
    echo "Failed to unzip AWS CLI v2" >&2
    exit 1
  fi
else
  echo "Failed to download AWS CLI v2" >&2
  exit 1
fi

# Fetch and install the binary from S3
echo "Fetching and installing the binary from S3..."
if aws s3 cp s3://urlshopbins/BackBinary /usr/local/bin/BackBinary; then
  echo "Binary fetched and installed successfully."
  if sudo chmod +x /usr/local/bin/BackBinary; then
    echo "Binary made executable."
  else
    echo "Failed to make binary executable" >&2
    exit 1
  fi
else
  echo "Failed to fetch binary from S3" >&2
  exit 1
fi

# Configure Nginx
echo "Configuring Nginx..."
if echo 'server { listen 80; location / { proxy_pass http://localhost:8080; } }' | sudo tee /etc/nginx/sites-available/default; then
  echo "Nginx configuration updated."
  if sudo systemctl restart nginx; then
    echo "Nginx restarted successfully."
  else
    echo "Failed to restart Nginx" >&2
    exit 1
  fi
else
  echo "Failed to configure Nginx" >&2
  exit 1
  fi

# Set up cron jobs to ensure the binary is always fetched on reboot
echo "Setting up cron jobs..."
(crontab -l 2>/dev/null; echo '@reboot aws s3 cp s3://urlshopbins/BackBinary /usr/local/bin/BackBinary') | crontab - || { echo 'Failed to set up cron job to fetch binary'; exit 1; }
(crontab -l 2>/dev/null; echo '@reboot sudo chmod +x /usr/local/bin/BackBinary') | crontab - || { echo 'Failed to set up cron job to make binary executable'; exit 1; }
(crontab -l 2>/dev/null; echo '@reboot sudo systemctl restart BackBinary.service') | crontab - || { echo 'Failed to set up cron job to restart service'; exit 1; }



# Execute BackBinary
echo "Executing BackBinary..."
sudo /usr/local/bin/BackBinary &
if [ $? -eq 0 ]; then
  echo "BackBinary started successfully."
else
  echo "Failed to start BackBinary" >&2
  exit 1
fi
	`

	instance.AddUserData(jsii.String(userData))

	return stack
}

func main() {
	app := awscdk.NewApp(nil)
	UrlShopStack(app, "UrlShopStack", &InfrastructureStackProps{
		StackProps: awscdk.StackProps{
			Env: &awscdk.Environment{
				Region:  jsii.String("ap-southeast-1"),
				Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
			},
		},
	})
	app.Synth(nil)
}
