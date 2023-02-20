package handlers

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"go-usermgmt-grpc/db/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

type handler struct {
	DB      *gorm.DB
	Account models.Accounts
}

func NewDB(db *gorm.DB) handler {
	return handler{
		DB: db}
}

func (h *handler) CreateUser(account *models.Accounts) error {

	err := account.BeforeSave()
	if err != nil {
		log.Fatal("connot creat hash password")
	}
	if result := h.DB.Create(&account); result.Error != nil {
		var perr *pgconn.PgError
		errors.As(result.Error, &perr)
		switch perr.Code {
		case "23505":
			log.Println(perr.Detail)
		default:
			log.Fatal("cannot create user")
		}
	}
	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (h *handler) LoginCheck(username string, password string) (bool, error) {
	var err error
	user := h.Account
	err = h.DB.Model(h.Account).Where("username = ?", username).Take(&user).Error
	if err != nil {
		log.Fatal("user is not exsist: ", err)
	}
	err = VerifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return false, err
	}
	return true, nil
}
