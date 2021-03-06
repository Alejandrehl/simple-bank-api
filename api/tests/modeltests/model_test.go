package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/alejandrehl/simple-bank-api/api/controllers"
	"github.com/alejandrehl/simple-bank-api/api/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var userInstance = models.User{}
var accountInstance = models.Account{}
var entryInstance = models.Entry{}
var transferInstance = models.Transfer{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {
	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
	server.DB, err = gorm.Open(TestDbDriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", TestDbDriver)
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {
	refreshUserTable()

	user := models.User{
		Nickname: "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return user, nil
}

func seedUsers() error {
	users := []models.User{
		{
			Nickname: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		{
			Nickname: "Kenny Morris",
			Email:    "kenny@gmail.com",
			Password: "password",
		},
	}

	for i := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func refreshUserAndAccountTable() error {
	err := server.DB.DropTableIfExists(&models.User{}, &models.Account{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Account{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOneAccount() (models.Account, error) {

	err := refreshUserAndAccountTable()
	if err != nil {
		return models.Account{}, err
	}
	user := models.User{
		Nickname: "Sam Phil",
		Email:    "sam@gmail.com",
		Password: "password",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Account{}, err
	}
	account := models.Account{
		Name:    "Account name",
		Balance:  1000,
		OwnerID: user.ID,
	}
	err = server.DB.Model(&models.Account{}).Create(&account).Error
	if err != nil {
		return models.Account{}, err
	}
	return account, nil
}

func seedUsersAndAccounts() ([]models.User, []models.Account, error) {
	var err error

	if err != nil {
		return []models.User{}, []models.Account{}, err
	}
	var users = []models.User{
		{
			Nickname: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		{
			Nickname: "Magu Frank",
			Email:    "magu@gmail.com",
			Password: "password",
		},
	}
	var accounts = []models.Account{
		{
			Name:    "Account 1",
			Balance:  1000,
			OwnerID: users[0].ID,
		},
		{
			Name:    "Account 2",
			Balance:  1000,
			OwnerID: users[1].ID,
		},
	}

	for i := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		accounts[i].OwnerID = users[i].ID

		err = server.DB.Model(&models.Account{}).Create(&accounts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed accounts table: %v", err)
		}
	}
	return users, accounts, nil
}