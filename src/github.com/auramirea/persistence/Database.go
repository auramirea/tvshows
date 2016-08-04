package persistence

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"sync"
	"log"
	"time"
)
type userRepository struct {}
var instance *userRepository
var db *gorm.DB

var once sync.Once
func GetUserRepository() *userRepository {
	once.Do(func() {
		instance = &userRepository{}
		db = openConnection()
	})
	return instance
}

type userRepositoryInterface interface {
	CreateUser(UserEntity) UserEntity
	DeleteUser(uint)
	FindAllUsers() []UserEntity
	FindUser() UserEntity
}

type UserEntity struct {
	ID        	uint `gorm:"primary_key"`
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
	FirstName       string  `gorm:"size:255"` // Default size for string is 255, reset it with this tag
	LastName        string  `gorm:"size:255"` // Default size for string is 255, reset it with this tag
	Email           string  `gorm:"size:255"`
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
	db.Create(&user)
	return user
}

func (*userRepository) DeleteUser(userId string) {
	user := UserEntity{}
	if db.First(&user, userId).Error != nil {
		log.Println("Couldn't find user with id", userId)
		return
	}
	if err := db.Delete(&user).Error; err != nil {
		log.Println("Couldn't delete user with id", userId)
		return
	}
}

func (*userRepository) FindUser(userId string) *UserEntity {
	user := UserEntity{}
	if db.First(&user, userId).Error != nil {
		log.Println("Couldn't find user with id", userId)
		return nil
	}
	return &user
}


func (*userRepository) FindAllUsers() []UserEntity {
	users := []UserEntity{}
	db.Find(&users)
	return users
}

