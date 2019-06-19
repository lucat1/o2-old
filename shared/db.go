package shared

import (
	"log"

	"github.com/jinzhu/gorm"
	// SQLite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

// OpenDatabase opens the gorm database
func OpenDatabase() {
	d, err := gorm.Open("sqlite3", "db.db")
	if err != nil {
		log.Fatal(err)
	}

	db = d
	db.LogMode(true)
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Repository{})
}

// GetDatabase returns the database
func GetDatabase() *gorm.DB {
	return db
}
