package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Terracode-Dev/terraui-back/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// var jwtSecret = "N9dnx3hLakwCvns5hY0aEjihuBqtALpBDahXyRRMiS4="

//TODO: [Client Siode] -- For both http responmses like unauth and internal server err, check for unauth status (401) or serer error (500) and redirect to login babaa...

func AddAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("inside Auth middleware")
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
			return []byte(os.Getenv("JKEY")), nil
		})

		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
		}

		userID := (*claims)["id"].(string)
		tenantID := (*claims)["tid"].(string)
		//fmt.Println("User ID in auth middleware:", userID) //TODO: remove this, its success...
		userData, err := database.FetchUserData(userID, tenantID) //TODO: change this to fetchUserData after TESTING
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user data")
		}

		// Set user data in context
		c.Set("user", userData)

		return next(c)
	}
}
