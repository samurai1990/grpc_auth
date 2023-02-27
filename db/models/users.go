package models

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type Users struct {
	Id        uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Username  string                `gorm:"unique;type:varchar(255);not null" json:"username"`
	Password  string                `gorm:"not null;type:varchar(500)" json:"password"`
	Email     string                `gorm:"uniqueIndex;type:varchar(255);not null" json:"email"`
	IsActive  bool                  `gorm:"default:true;type:boolean" json:"is_active"`
	IsAdmin   bool                  `gorm:"default:false;type:boolean" json:"is_admin"`
	CreatedAt time.Time             `gorm:"default:current_timestamp" json:"created"`
	UpdatedAt time.Time             `gorm:"default:current_timestamp" json:"updated"`
	DeletedAt time.Time             `json:"deleted"`
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag;DeletedAtField:DeletedAt"`
}

func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&Users{})
	return err
}

func NewUser(username string, password string, email string, isAdmin bool, isActive bool) (*Users, error) {
	user := &Users{
		Username: username,
		Password: password,
		Email:    email,
		IsAdmin:  isAdmin,
		IsActive: isActive,
	}
	return user, nil
}

func encrtyptPasswords(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	hashedPassword := hex.EncodeToString(h.Sum(nil))
	return hashedPassword
}

func (user *Users) BeforeCreate(tx *gorm.DB) (err error) {
	user.Username = strings.TrimSpace(user.Username)
	user.Password = encrtyptPasswords(user.Password)
	return nil
}

func (user *Users) IsCorrectPassword(password string) bool {
	return encrtyptPasswords(password) == user.Password
}

func (user *Users) BeforeDelete(tx *gorm.DB) (err error) {
	//impelement is_deleted
	if user.IsAdmin {
		return errors.New("admin user not allowed to delete")
	}
	return
}
