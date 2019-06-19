package shared

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// User represents a user in the database
type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`

	Username  string `gorm:"unique_index"`
	Firstname string
	Lastname  string
	Email     string `gorm:"unique"`
	Password  string `gorm:"type:varchar(100)"`

	// Profile specific stuff
	Description  string
	Repositories []Repository
}

// BeforeCreate is used to generate the user's uuid
func (*User) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4())
}
