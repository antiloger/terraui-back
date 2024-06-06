package api

import (
	"github.com/Terracode-Dev/terraui-back/api/handlers"
)

func (s *Server) Router_v1() {
	s.Echo.GET("/hi", handlers.Hello)
}
