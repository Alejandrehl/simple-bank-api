package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alejandrehl/simple-bank-api/api/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Println("Cannot connect to postgres database")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Println("We are connected to the postgres database")
	}

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Account{}, &models.Entry{}, &models.Transfer{}) //database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}