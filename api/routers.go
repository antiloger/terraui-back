package api

import "github.com/Terracode-Dev/terraui-back/api/middleware"

func (s *Server) Router_v1() {
	s.Echo.GET("/hi", Hello)

	s.Echo.POST("/login", s.UserLogin)
	s.Echo.GET("/getalltables", s.GetAllTables, middleware.AddAuth)
}
