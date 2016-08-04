package persistence

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"sync"
	"log"
	"time"
	"github.com/auramirea/service"
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

type UserEntity struct {
	ID        	uint `gorm:"primary_key"`
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
	FirstName       string  `gorm:"size:255"` // Default size for string is 255, reset it with this tag
	LastName        string  `gorm:"size:255"` // Default size for string is 255, reset it with this tag
	Email           string  `gorm:"size:255"`
	Shows 		[]Show  `gorm:"ForeignKey:UserId"`
}
type Show struct {
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
	Name            string  `gorm:"size:255"` // Default size for string is 255, reset it with this tag
	ExternalId      int `gorm:"primary_key"`
	UserId          uint `gorm:"primary_key"`
}
func (UserEntity) TableName() string {
	return "appuser"
}
func (Show) TableName() string {
	return "show"
}

func openConnection() (*gorm.DB) {
	db, err := gorm.Open("postgres", DB_URL)
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&UserEntity{}, &Show{})

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
	db.Model(&user).Association("Shows").Find(&user.Shows)
	return &user
}


func (*userRepository) FindAllUsers() []UserEntity {
	users := []UserEntity{}
	result := []UserEntity{}
	db.Find(&users)
	for _, user:= range(users) {
		db.Model(&user).Association("Shows").Find(&user.Shows)
		result = append(result, user)
	}
	return result
}

func (u *userRepository) AddShow(userId string, show *service.Show) *UserEntity {
	user := UserEntity{}
	if db.First(&user, userId).Error != nil {
		log.Println("Couldn't find user with id", userId)
		return nil
	}
	showToSave := Show{ExternalId: show.Id, Name: show.Name, UserId: user.ID}
	db.Create(&showToSave)

	return u.FindUser(userId)
}

func (u *userRepository) DeleteShow(userId string, showId string) *UserEntity {
	user := u.FindUser(userId)
	if user == nil {
		return nil
	}

	show := Show{}
	db.First(&show, showId)
	db.Delete(&show)

	return u.FindUser(userId)
}

