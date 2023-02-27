package db_test

import (
	"go-usermgmt-grpc/db"
	"go-usermgmt-grpc/db/models"
	"go-usermgmt-grpc/service"
	"testing"

	"github.com/stretchr/testify/require"
)

func MockNewUser() *models.Users {
	return &models.Users{
		Username: "test",
		Password: "test",
		Email:    "test@tes.test",
		IsActive: false,
		IsAdmin:  false,
	}
}

func TestFindUser(t *testing.T) {
	sampleUser := MockNewUser()
	userStore := db.UserStore{}
	newUser, err := service.CreateUser(
		&userStore,
		sampleUser.Username,
		sampleUser.Password,
		sampleUser.Email,
		sampleUser.IsActive,
		sampleUser.IsAdmin,
	)
	require.Nil(t, err)
	require.NotEmpty(t, newUser)

	query, err := userStore.Find(sampleUser.Username)
	require.Nil(t, err)
	require.Equal(t, query.Username, sampleUser.Username)

	err = userStore.Delete(query)
	require.Nil(t, err)
}
