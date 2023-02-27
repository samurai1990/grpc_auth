package service

import (
	"go-usermgmt-grpc/db"
	"go-usermgmt-grpc/db/models"
)

func CreateUser(userStore db.UserStoreInterface, username string, password string, email string, isAdmin bool, isActive bool) (*models.Users, error) {
	newUser, err := models.NewUser(username, password, email, isAdmin, isActive)
	if err != nil {
		return nil, err
	}
	user, err := userStore.Save(newUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func FindUser(serStore db.UserStoreInterface, username string) *models.Users {
	user, err := serStore.Find(username)
	if err != nil {
		return nil
	}
	return user
}

// func (user *User) Clone() *User {
// 	return &User{
// 		Username:       user.Username,
// 		HashedPassword: user.HashedPassword,
// 		Role:           user.Role,
// 	}
// }
