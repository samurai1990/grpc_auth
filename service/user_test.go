package service_test

import (
	"go-usermgmt-grpc/db"
	"go-usermgmt-grpc/db/models"
	"go-usermgmt-grpc/service"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func MockNewUser() *models.Users {
	return &models.Users{
		Username: "test",
		Password: "test",
		Email:    "test@tes.test",
		IsActive: false,
		IsAdmin:  true,
	}
}

func TestCreateUser(t *testing.T) {

	sampleUser := MockNewUser()
	userStore := &db.UserStore{}
	newUser, err := service.CreateUser(
		userStore,
		sampleUser.Username,
		sampleUser.Password,
		sampleUser.Email,
		sampleUser.IsActive,
		sampleUser.IsAdmin,
	)
	log.Println(sampleUser.Password)
	require.Nil(t, err)
	require.NotNil(t, newUser)
}