package api

import "github.com/Terracode-Dev/terraui-back/api/middleware"

func (s *Server) Router_v1() {
	s.Echo.GET("/hi", Hello)

	s.Echo.POST("/login", s.UserLogin)
	s.Echo.POST("/register", s.UserRegister)
	s.Echo.Use(middleware.AddAuth)
	s.Echo.GET("/tabledata", s.GetAllItems)
	s.Echo.GET("/getalltables", s.GetAllTables)
	s.Echo.GET("/itemdata", s.GetItem)
	s.Echo.GET("/userdata", s.GetUserData)
	s.Echo.GET("/paymentdata", s.PaymentData)
	s.Echo.GET("/orderdata", s.OrderData)
	s.Echo.POST("/newlink", s.LinkGenarator)
	s.Echo.DELETE("/deletetable", s.DeleteTable)
	s.Echo.DELETE("/deleteitem", s.DeleteItem)
	s.Echo.PUT("/updatetable", s.UpdateTable)
	s.Echo.PUT("/updateitem", s.UpdateItem)
}
