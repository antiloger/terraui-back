package api

import (
	"github.com/Terracode-Dev/terraui-back/api/handlers"
	"github.com/Terracode-Dev/terraui-back/api/middleware"
)

func (s *Server) Router_v1() {
	s.Echo.GET("/hi", handlers.Hello)

	//TODO: Remove this (fir Auth test only))

	Grp := s.Echo.Group("/authtest")
	Grp.GET("/", handlers.TestAuth, middleware.AddAuth) // Add : middleware for auth added here
}
