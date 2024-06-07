package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type UserData struct {
	UID   string `json:"id"`
	uname string `json:"name"`
	email string `json:"email"`
}

// fetchUserData fetches user data based on user ID
func fetchUserData(userID string) (*UserData, error) {

	var userData UserData
	// Implement your user data fetching logic here
	// This is a mock implementation
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %w", err)
	}

	// Create a DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	// Set up the query input parameters
	params := &dynamodb.QueryInput{
		TableName:              aws.String("YourTableName"),
		KeyConditionExpression: aws.String("user_id = :user_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_id": &types.AttributeValueMemberS{Value: userID},
		},
	}

	// Execute the query
	result, err := svc.Query(context.TODO(), params)
	if err != nil {
		return nil, fmt.Errorf("failed to query DynamoDB, %w", err)
	}

	// Check if we got any items
	if len(result.Items) == 0 {
		return nil, fmt.Errorf("no items found")
	}

	// Deserialize the first item (assuming the user ID is unique)
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &userData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data, %w", err)
	}

	return &userData, nil
}

var jwtSecret = []byte("your-secret-key")

func AddAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid Authorization header")
		}

		tokenString := authHeader[len("Bearer "):]

		claims := &jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
		}

		// Fetch user data based on claims (for example, using the user ID from the token claims)
		userID := (*claims)["id"].(string)
		userData, err := fetchUserData(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user data")
		}

		// Set user data in context
		c.Set("user", userData)

		return next(c)
	}
}
