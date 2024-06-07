package middleware

import (
	"errors"
	"net/http"

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
	// Implement your user data fetching logic here
	// This is a mock implementation
	return &UserData{
		UID:   userID,
		uname: "John Doe",
		email: "johndoe@example.com",
	}, nil
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
