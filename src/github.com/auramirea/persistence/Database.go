package persistence

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"sync"
)
type userRepository struct {}
var instance *userRepository

var once sync.Once
func GetUserRepository() *userRepository {
	once.Do(func() {
		instance = &userRepository{}
	})
	return instance
}

type userRepositoryInterface interface {
	CreateUser(UserEntity) UserEntity
}

type UserEntity struct {
	gorm.Model
	FirstName        string  `gorm:"size:255"` // Default size for string is 255, reset it with this tag
	LastName         string  `gorm:"size:255"` // Default size for string is 255, reset it with this tag
	Email            string  `gorm:"size:255"`
}
func (UserEntity) TableName() string {
	return "appuser"
}

func openConnection() (*gorm.DB) {
	db, err := gorm.Open("postgres", DB_URL)
	if err != nil {
		panic("failed to connect database")
	}
	return db

}
func (*userRepository) CreateUser(user UserEntity) UserEntity {
	db := openConnection()
	db.Create(&user)
	defer db.Close()
	return user
}
