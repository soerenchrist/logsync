package model

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"time"
)

type OperationType string

const (
	Deleted  OperationType = "D"
	Created  OperationType = "C"
	Modified OperationType = "M"
)

type ChangeLogEntry struct {
	GraphName string    `gorm:"primaryKey"`
	FileId    string    `gorm:"primaryKey"`
	Timestamp time.Time `gorm:"primaryKey"`
	FileName  string
	Operation OperationType
}

func CreateDb(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&ChangeLogEntry{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
