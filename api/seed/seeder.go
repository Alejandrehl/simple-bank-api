package seed

import (
	"log"

	"github.com/alejandrehl/simple-bank-api/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	{
		Nickname: "alejandrehl",
		Email:    "alejandrehl@icloud.com",
		Password: "asd123.",
	},
	{
		Nickname: "alehernandezdev",
		Email:    "alehernandezdev@gmail.com",
		Password: "asd123.",
	},
}

func Load(db *gorm.DB) {
	for i := range users {
		var err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}