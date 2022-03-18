package controllers

import "github.com/alejandrehl/simple-bank-api/api/middlewares"

func (s *Server) initializeRoutes() {
	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Accounts routes
	s.Router.HandleFunc("/accounts", middlewares.SetMiddlewareJSON(s.Create)).Methods("POST")
	s.Router.HandleFunc("/accounts", middlewares.SetMiddlewareJSON(s.GetAll)).Methods("GET")
	s.Router.HandleFunc("/accounts/{id}", middlewares.SetMiddlewareJSON(s.GetById)).Methods("GET")
	s.Router.HandleFunc("/accounts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.Update))).Methods("PUT")
	s.Router.HandleFunc("/accounts/{id}", middlewares.SetMiddlewareAuthentication(s.Delete)).Methods("DELETE")
}