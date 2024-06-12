package api

import (
	"github.com/Terracode-Dev/terraui-back/database"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Addr string
	Echo *echo.Echo
	DB   *database.DB
}

func NewServer(addr string, db *database.DB) Server {
	return Server{
		Addr: addr,
		Echo: echo.New(),
		DB:   db,
	}
}

func (s *Server) Run() {
	s.Router_v1() // init routers
	s.Echo.Logger.Fatal(s.Echo.Start(s.Addr))
}
