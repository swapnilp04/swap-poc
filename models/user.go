package models

import (
	"encoding/hex"
	"fmt"
	"swapnil-ex/models/db"
	"swapnil-ex/swapErr"
	
	"gorm.io/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/scrypt"
	"time"
)

type User struct {
	ID              int    `json:"id"`
	Username        string `json:"username"`
	Salt            string `json:"-"`
	Password        string `json:"-"`
	ConfirmPassword string `json:"-" gorm:"-"`
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
  DeletedAt 			gorm.DeletedAt `gorm:"index"`
}

func migrateUser() {
	fmt.Println("migrating user..")
	err := db.Driver.AutoMigrate(&User{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func (u *User) FindUserByUsername(username string) error {
	err := db.Driver.First(u, "username = ?", username).Error
	return err
}

func (u *User) ValidPassword(password string) error {
	hash, err := scrypt.Key([]byte(password), []byte(u.Salt), 32768, 8, 1, 32)
	if err != nil {
		return errors.Wrap(err, "ValidPassword(scrypt.Key)")
	}

	if hex.EncodeToString(hash) != u.Password {
		return swapErr.ErrInvalidUser
	}

	return nil
}

func (u *User) Save() error {
	err := db.Driver.Save(u).Error
	return err
}

func (u *User) Validate() error {
	if u.ConfirmPassword != u.Password {
		return swapErr.ErrPasswordMisMatch
	}

	////Other check///
	return nil
}

func (u *User) Load() error {
	err := db.Driver.Find(u, "id = ?", u.ID).Error
	return err
}
