package database

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	apitypes "github.com/Terracode-Dev/terraui-back/types"
	apierror "github.com/Terracode-Dev/terraui-back/util/apierrors"
)

type DB struct {
	Conn *dynamodb.Client
}

func StartDB() (*DB, error) {
	client, err := InitDynamoDBClient()
	if err != nil {
		return nil, err
	}
	return &DB{
		Conn: client,
	}, nil
}

// ---- dynamoDB Implementations----

func InitDynamoDBClient() (*dynamodb.Client, error) {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}

// ADDED the USER FETCH METHOD and user struct

func (db *DB) CheckUser(user *apitypes.UserLogin) (*apitypes.User, error) {
	iquery := &dynamodb.QueryInput{
		TableName:              aws.String("user_details"),
		KeyConditionExpression: aws.String("user_id = :e AND tenant_id = :t"),
		FilterExpression:       aws.String("userkey = :p"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":e": &types.AttributeValueMemberS{Value: user.Userid},
			":t": &types.AttributeValueMemberS{Value: user.Tenant},
			":p": &types.AttributeValueMemberS{Value: user.Password}, // TODO: add hash
		},
		ProjectionExpression: aws.String("user_id, #role, #subscription, #useremail, #username"),
		ExpressionAttributeNames: map[string]string{
			"#role":         "role",
			"#subscription": "subscription",
			"#useremail":    "useremail",
			"#username":     "username",
		},
	}

	result, err := db.Conn.Query(context.TODO(), iquery)
	if err != nil {
		return nil, err
	}

	fmt.Println(result)
	u := new(apitypes.User)

	if len(result.Items) == 0 {
		return nil, apierror.ErrAuthFail
	}

	fmt.Println("lol")
	err = attributevalue.UnmarshalMap(result.Items[0], u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// func (db *DB) FetchUserData(userID string, tenantID string) (*UserData, error) {
// 	var userData UserData
//
// 	// Set up the query input parameters
// 	SQRY := &dynamodb.QueryInput{
// 		TableName:              aws.String("user_details"),
// 		KeyConditionExpression: aws.String("user_id = :user_id AND tenant_id = :tenant_id"),
// 		ExpressionAttributeValues: map[string]types.AttributeValue{
// 			":user_id":   &types.AttributeValueMemberS{Value: userID},
// 			":tenant_id": &types.AttributeValueMemberS{Value: *aws.String(tenantID)},
// 		},
// 	}
//
// 	// Execute the query
// 	result, err := db.Conn.Query(context.TODO(), SQRY)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query DynamoDB, %w", err)
// 	}
//
// 	// Check if we got any items
// 	if len(result.Items) == 0 {
// 		return nil, fmt.Errorf("no items found")
// 	}
//
// 	// Deserialize the first item (assuming the user ID is unique)
// 	err = attributevalue.UnmarshalMap(result.Items[0], &userData)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal data, %w", err)
// 	}
//
// 	return &userData, nil
// }
