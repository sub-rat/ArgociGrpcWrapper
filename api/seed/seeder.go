package seed
import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/sub-rat/ArogciGrpcWrapper/api/models"
)

var users = []models.User{
	models.User{
		Nickname: "admin",
		Email:    "admin@gmail.com",
		Password: "password",
	},
	models.User{
		Nickname: "admin1",
		Email:    "admin1@gmail.com",
		Password: "password",
	},
}

func Load(db *gorm.DB) {

	for  i , _  :=  range  users {
		err := db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	}
}