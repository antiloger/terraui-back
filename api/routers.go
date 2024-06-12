package api

func (s *Server) Router_v1() {
	s.Echo.GET("/hi", Hello)

	s.Echo.POST("/login", s.UserLogin)
}
