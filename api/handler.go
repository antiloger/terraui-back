package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Terracode-Dev/terraui-back/types"
	"github.com/Terracode-Dev/terraui-back/util"
	apierror "github.com/Terracode-Dev/terraui-back/util/apierrors"
	"github.com/labstack/echo/v4"
)

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello TerraCode!")
}

func (s *Server) UserLogin(c echo.Context) error {
	ul := new(types.UserLogin)
	if err := c.Bind(ul); err != nil {
		return err
	}
	fmt.Println("test")
	res, err := s.DB.CheckUser(ul)
	if err != nil {
		if errors.Is(err, apierror.ErrAuthFail) {
			return c.JSON(http.StatusUnauthorized, types.NewResErr("Invalid password or user email", 9))
		}
		return c.JSON(http.StatusInternalServerError, types.NewResErr("Internal server issue", 91))
	}

	//---set cookies---
	cookie := new(http.Cookie)
	cookie.Name = "terralinkck"
	jk, err := util.GetToken(res)
	if err != nil {
		return nil
	}
	cookie.Value = jk
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
	return c.JSON(http.StatusAccepted, types.NewResErr("", 0))
}

// : remove this ( created only for auth testing )
// func TestAuth(c echo.Context) error {
// 	user, ok := c.Get("user").(*database.UserData)
// 	// user, ok := userData.(*database.UserData)
// 	if !ok {
// 		// Handle the case where the type assertion fails
// 		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user data")
// 	}
// 	return c.String(http.StatusOK, "Hello, "+user.Uname+"!")
// }
