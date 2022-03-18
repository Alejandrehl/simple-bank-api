package controllers

import (
	"net/http"

	"github.com/alejandrehl/simple-bank-api/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to the Simple Bank API")
}