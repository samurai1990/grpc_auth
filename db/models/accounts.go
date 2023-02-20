package models

import (
	"strings"
	"time"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Accounts struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Username  string    `gorm:"unique;type:varchar(255);not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updated"`
	IsActive  bool      `gorm:"not null" json:"is_active"`
	IsDeleted bool      `gorm:"not null" json:"is_deleted"`
	IsAdmin   bool      `gorm:"not null" json:"is_admin"`
}

func (account *Accounts) BeforeSave() error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.Username = strings.TrimSpace(account.Username)
	account.Password = string(hashedPassword)
	return nil
}
