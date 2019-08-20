package db

import "time"

// Media represents a Media item in the database.
type Media struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Length uint64 `gorm:"not null"`
	Title  string `gorm:"not null"`
	Type   string `gorm:"primary_key"`
}
