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
	GraphName   string        `gorm:"primaryKey" json:"graph_name"`
	FileId      string        `gorm:"primaryKey" json:"file_id"`
	Timestamp   time.Time     `gorm:"primaryKey" json:"timestamp"`
	Transaction string        `json:"transaction"`
	Operation   OperationType `json:"operation"`
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
