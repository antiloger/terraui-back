package api

import (
	"errors"
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
	res, err := s.DB.CheckUser(ul)
	if err != nil {
		if errors.Is(err, apierror.ErrAuthFail) {
			return c.JSON(http.StatusUnauthorized, types.NewResErr("Invalid user email", 98))
		}
		return c.JSON(http.StatusInternalServerError, types.NewResErr("Internal server issue", 91))
	}

	if !(util.HashCheck(ul.Password, res.Userkey)) {
		return c.JSON(http.StatusUnauthorized, types.NewResErr("Invalid password", 9))
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
	return c.JSON(http.StatusAccepted, types.NewResErr("StatusAccepted", 12))
}

func (s *Server) GetAllTables(c echo.Context) error {
	user, ok := c.Get("Auth").(*types.AuthUser)
	if !ok {
		return c.JSON(http.StatusUnauthorized, types.NewResErr("Authantication Failed", 91))
	}
	set, err := s.DB.GetAllUserTables(user.Userid)
	if err != nil {
		if errors.Is(err, apierror.ErrZeroData) {
			return c.JSON(http.StatusAccepted, types.NewResErr("No Table found for this user", 21))
		}
		return err
	}

	return c.JSON(http.StatusOK, types.NewUserTables(user.Userid, set).Format())
}

func (s *Server) GetAllItems(c echo.Context) error {
	t := new(types.GetAllItemRes)
	if err := c.Bind(t); err != nil {
		return err
	}
	user, ok := c.Get("Auth").(*types.AuthUser)
	if !ok {
		return c.JSON(http.StatusUnauthorized, types.NewResErr("Authantication Failed", 92))
	}

	tableinfo, err := s.DB.CheckItems(user.Userid, t.Table_id)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.NewResErr("Authantication Failed", 91))
	}

	tabledata, err := s.DB.GetAllItem(t.Table_id)
	if err != nil {
		return c.JSON(http.StatusAccepted, types.NewResErr("empty", 12))
	}

	tableitem := types.NewTableItems(tableinfo, tabledata)

	return c.JSON(http.StatusOK, tableitem.Format("table item", 0))
}

func (s *Server) UserRegister(c echo.Context) error {
	u := new(types.UserRegister)
	if err := c.Bind(u); err != nil {
		return err
	}
	u.Userid = util.NewUUID()
	u.Role = "user"
	key, err := util.HashKey(u.Userkey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.NewResErr("Internal Server error", 99))
	}
	u.Userkey = key

	err = s.DB.CreateUser(u)
	if err != nil {
		return nil
	}

	return c.JSON(http.StatusCreated, types.NewResErr("user created", 11))
}

func (s *Server) GetItem(c echo.Context) error
func (s *Server) GetUserData(c echo.Context) error
func (s *Server) LinkGenarator(c echo.Context) error
func (s *Server) DeleteTable(c echo.Context) error
func (s *Server) DeleteItem(c echo.Context) error
func (s *Server) UpdateTable(c echo.Context) error
func (s *Server) UpdateItem(c echo.Context) error
func (s *Server) PaymentData(c echo.Context) error
func (s *Server) OrderData(c echo.Context) error

// : remove this ( created only for auth testing )
// func TestAuth(c echo.Context) error {
// 	user, ok := c.Get("user").(*database.UserData)
// 	// user, ok := userData.(*database.UserData)
// 	if !ok {
// 		// Handle the case where the type assertion fails
// 		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user data")
// 	}
// 	return c.String(http.StatusOK, "Hello, "+user.Uname+"!")
//
