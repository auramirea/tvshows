package service

import (
	"sync"
	//"github.com/auramirea/persistence"
)

type userService struct {}
var instance *userService

var once sync.Once
func GetUserService() *userService {
	once.Do(func() {
		instance = &userService{}
	})
	return instance
}

//type userServiceInterface interface {
//	CreateUser(User) User
//	DeleteUser(string)
//	FindAllUsers() []User
//	FindUser() User
//}
//
//type User struct {
//	FirstName string `json:"firstName"`
//	LastName string `json:"lastName"`
//	Email string `json:"email"`
//	Id uint `json:"id"`
//}
//
//func (serviceInstance *userService) CreateUser(user User) User {
//	repo := persistence.GetUserRepository()
//	return convertToUser(repo.CreateUser(convertToUserEntity(user)))
//}
//
//func (serviceInstance *userService) DeleteUser(userId string) {
//	persistence.GetUserRepository().DeleteUser(userId)
//}
//
//func (serviceInstance *userService) FindAllUsers() []User {
//	userEntities := persistence.GetUserRepository().FindAllUsers()
//	users := make([]User, len(userEntities))
//	for i, userEntity := range userEntities {
//		users[i] = convertToUser(userEntity)
//	}
//	return users
//}
//
//func (serviceInstance *userService) FindUser(userId string) User {
//	return convertToUser(persistence.GetUserRepository().FindUser(userId))
//}
//
//func convertToUserEntity(user User) persistence.UserEntity {
//	return persistence.UserEntity{FirstName: user.FirstName, LastName: user.LastName, Email: user.Email}
//}
//func convertToUser(user persistence.UserEntity) User {
//	return User{FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Id: user.ID}
//}