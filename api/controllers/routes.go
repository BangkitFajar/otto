package controllers

import "tes/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Register
	s.Router.HandleFunc("/register", middlewares.SetMiddlewareJSON(s.Register)).Methods("POST")
	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/user", middlewares.SetMiddlewareAuthentication(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/historys", middlewares.SetMiddlewareAuthentication(s.GetHistoryTransaction)).Methods("GET")

	//Balance routes
	s.Router.HandleFunc("/generate-va", middlewares.SetMiddlewareAuthentication(s.GenerateVA)).Methods("POST")
	s.Router.HandleFunc("/checking", middlewares.SetMiddlewareAuthentication(s.Checking)).Methods("POST")
}
