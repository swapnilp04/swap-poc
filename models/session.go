package models

import (
	"fmt"
	"swapnil-ex/models/db"
	"time"
)

type Session struct {
	ID        string
	UserID    int
	CreatedAt time.Time
	ExpiresAt time.Time
}

func migrateSession() {
	fmt.Println("migrating session..")

	err := db.Driver.AutoMigrate(&Session{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func (s *Session) Save() error {
	err := db.Driver.Save(s).Error
	return err
}

func (s *Session) Load() error {
	err := db.Driver.Find(s, "id = ?", s.ID).Error
	return err
}

func (s *Session) Valid() bool {
	return s.ExpiresAt.After(time.Now())
}

func (s *Session) Destroy() error {
	err := db.Driver.Delete(s).Error
	return err
}
