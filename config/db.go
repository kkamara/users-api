package config

import (
	"github.com/kkamara/users-api/schemas/userSchema"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenDB() (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	db.AutoMigrate(&userSchema.UserSchema{})
	return
}
