package handlers

import (
	"net/http"

	"github.com/Terracode-Dev/terraui-back/database"
	"github.com/labstack/echo/v4"
)

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello TerraCode!")
}

// TODO: remove this ( created only for auth testing )
func TestAuth(c echo.Context) error {
	user, ok := c.Get("user").(*database.UserData)
	//user, ok := userData.(*database.UserData)
	if !ok {
		// Handle the case where the type assertion fails
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user data")
	}
	return c.String(http.StatusOK, "Hello, "+user.Uname+"!")
}
