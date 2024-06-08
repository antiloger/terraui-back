package database

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DB struct {
	ConnStr string
}

func NewDB(connstr string) DB {
	return DB{
		ConnStr: connstr,
	}
}

// add db run and return error
func (d *DB) Run() {
}

// ---- dynamoDB codes----
var svc *dynamodb.Client

func InitDynamoDBClient() error {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("unable to load SDK config, %w", err)
	}

	// Create a DynamoDB client
	svc = dynamodb.NewFromConfig(cfg)
	return nil

}

// ADDED the USER FETCH METHOD and user struct
// ----MODEL STRUCTS----
type UserData struct {
	UID   string `dynamodbav:"user_id"`
	Uname string `dynamodbav:"name"`
	Email string `dynamodbav:"email"`
}

// --------DATA METHODS--------
func FetchUserData(userID string, tenantID string) (*UserData, error) {

	var userData UserData

	// Set up the query input parameters
	SQRY := &dynamodb.QueryInput{
		TableName:              aws.String("user_details"),
		KeyConditionExpression: aws.String("user_id = :user_id AND tenant_id = :tenant_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id":   &types.AttributeValueMemberS{Value: userID},
			":tenant_id": &types.AttributeValueMemberS{Value: *aws.String(tenantID)},
		},
	}

	// Execute the query
	result, err := svc.Query(context.TODO(), SQRY)

	if err != nil {
		return nil, fmt.Errorf("failed to query DynamoDB, %w", err)
	}

	// Check if we got any items
	if len(result.Items) == 0 {

		return nil, fmt.Errorf("no items found")
	}

	// Deserialize the first item (assuming the user ID is unique)
	err = attributevalue.UnmarshalMap(result.Items[0], &userData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data, %w", err)
	}

	return &userData, nil
}
