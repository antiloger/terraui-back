package database

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/labstack/gommon/log"

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
		KeyConditionExpression: aws.String("tenant_id = :t"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":t": &types.AttributeValueMemberS{Value: user.Tenant},
		},
		ProjectionExpression: aws.String("user_id, #role, #subscription, #userkey, #tenant_id, #username"),
		ExpressionAttributeNames: map[string]string{
			"#role":         "role",
			"#subscription": "subscription",
			"#tenant_id":    "tenant_id",
			"#username":     "username",
			"#userkey":      "userkey",
		},
	}

	result, err := db.Conn.Query(context.TODO(), iquery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(result)
	u := new(apitypes.User)

	if len(result.Items) == 0 {
		return nil, apierror.ErrAuthFail
	}

	err = attributevalue.UnmarshalMap(result.Items[0], u)
	if err != nil {
		return nil, err
	}
	fmt.Println(u.Userkey)

	return u, nil
}

// to create user
func (db *DB) CreateUser(u *apitypes.UserRegister) error {
	user, err := attributevalue.MarshalMap(u)
	if err != nil {
		log.Printf("\n %s \n", err)
		return err
	}

	iquery := &dynamodb.PutItemInput{
		TableName: aws.String("user_details"),
		Item:      user,
	}

	output, err := db.Conn.PutItem(context.TODO(), iquery)
	if err != nil {
		fmt.Println("ttt")
		log.Printf("\n %s \n", err)
		return err
	}

	fmt.Println(output)
	return nil
}

func (db *DB) GetAllUserTables(userid string) (*[]apitypes.TableInfo, error) {
	iquery := &dynamodb.QueryInput{
		TableName:              aws.String("users_tables"),
		KeyConditionExpression: aws.String("user_id = :i"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":i": &types.AttributeValueMemberS{Value: userid},
		},
		// ProjectionExpression: aws.String("user_id, #tablename, #discription, #lastdate, #color, #columns"),
		// ExpressionAttributeNames: map[string]string{
		//   "#tablename": "tablename",
		//   "#discription": "discription",
		//   "#lastdate": "lastdate",
		//   "#color": "color",
		//   "#columns": "columns",
		// },
	}

	result, err := db.Conn.Query(context.TODO(), iquery)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, apierror.ErrZeroData
	}
	fmt.Println(result.Items)
	info := new([]apitypes.TableInfo)

	err = attributevalue.UnmarshalListOfMaps(result.Items, info)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return info, nil
}

func (db *DB) CheckItems(userid string, tableid string) (*apitypes.TableInfo, error) {
	fmt.Println(userid, tableid)
	iquery := &dynamodb.QueryInput{
		TableName:              aws.String("users_tables"),
		KeyConditionExpression: aws.String("user_id = :i AND table_id = :t"), // TODO: change this "user_id"
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":i": &types.AttributeValueMemberS{Value: userid},
			":t": &types.AttributeValueMemberS{Value: tableid},
		},
	}

	result, err := db.Conn.Query(context.TODO(), iquery)
	if err != nil {
		return nil, err
	}
	fmt.Println(result)

	if len(result.Items) == 0 {
		return nil, apierror.ErrZeroData
	}
	info := new(apitypes.TableInfo)
	err = attributevalue.UnmarshalMap(result.Items[0], info)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return info, nil
}

func (db *DB) GetAllItem(tableid string) (*[]map[string]any, error) {
	iquery := &dynamodb.QueryInput{
		TableName:              aws.String("items_store"),
		KeyConditionExpression: aws.String("user_id = :i"), // TODO: change this "user_id"
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":i": &types.AttributeValueMemberS{Value: tableid},
		},
	}

	result, err := db.Conn.Query(context.TODO(), iquery)
	if err != nil {
		return nil, err
	}

	info := new([]map[string]any)
	err = attributevalue.UnmarshalListOfMaps(result.Items, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}
