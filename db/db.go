package db

import (
	"errors"
	"fmt"
	"go-usermgmt-grpc/db/models"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserStore struct {
	*gorm.DB
	user *models.Users
}

func NewInDBUserStore() *UserStore {
	return &UserStore{}
}

type UserStoreInterface interface {
	Save(user *models.Users) (*models.Users, error)
	Find(username string) (*models.Users, error)
	Delete(user *models.Users) error
	List() ([]*models.Users, error)
}

var ErrAlreadyExists = errors.New("username already exists")

func (handle *UserStore) NewConnection() (*gorm.DB, error) {

	DB_Host := "192.168.10.15"
	DbUser := os.Getenv("POSTGRES_DEVELOP_DB_USERNAME")
	DbPassword := os.Getenv("POSTGRES_DEVELOP_DB_PASSWORD")
	DbName := os.Getenv("POSTGRES_DEVELOP_DB_NAME")
	DbPort := 5432
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC", DB_Host, DbUser, DbPassword, DbName, DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Database: %v", err)
	}
	err = models.MigrateUsers(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (handle *UserStore) Save(user *models.Users) (*models.Users, error) {
	newDB, err := handle.NewConnection()
	if err != nil {
		return nil, err
	}
	result := newDB.Create(&user)
	if result.Error != nil {
		var perr *pgconn.PgError
		errors.As(result.Error, &perr)
		switch perr.Code {
		case "23505":
			return nil, ErrAlreadyExists
		default:
			return nil, result.Error
		}
	}
	log.Printf("user created with info: %v", user)
	return user, nil
}

func (handle *UserStore) Find(username string) (*models.Users, error) {
	newDB, err := handle.NewConnection()
	if err != nil {
		return nil, err
	}
	user := handle.user
	st := newDB.Model(&user).Where("username = ?", username).Where("is_active = ?", true).Take(&user).Error
	if st != nil {
		return nil, st
	}
	return user, nil
}

func (handle *UserStore) Delete(user *models.Users) error {
	newDB, err := handle.NewConnection()
	if err != nil {
		return err
	}
	return newDB.Unscoped().Delete(&user).Error
}

func (handle *UserStore) List() ([]*models.Users, error) {
	newDB, err := handle.NewConnection()
	if err != nil {
		return nil, err
	}
	users := []*models.Users{}
	listUsers := newDB.Find(&users)
	if listUsers.Error != nil {
		return nil, listUsers.Error
	}
	return users, nil
}
