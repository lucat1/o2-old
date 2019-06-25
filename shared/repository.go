package shared

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// Repository represents a git repo in the database
type Repository struct {
	gorm.Model
	Owner       string
	Name        string
	Description string
	MainBranch  string `gorm:"default:'master'"`

	// Permissions systems
	Restricted string `gorm:"default:'repo:push,repo:settings'"` // Slice of restricted scopes
	Granted    string `gorm:"default:'[]'"`
}

type permission struct {
	For string // The username of who the permission is for
	Key string // The actual permisson scope
}

func (r *Repository) BeforeCreate(scope *gorm.Scope) error {
	data := "[{\"For\":\"" + r.Owner + "\",\"Key\":\"repo:settings\"}]"
	return scope.SetColumn("Granted", data)
}

func ParsePermissions(r *Repository) *[]permission {
	var v []permission
	err := json.Unmarshal([]byte(r.Granted), &v)
	if err != nil {
		GetLogger().Warn(
			"Could not decode perimssions",
			zap.String("repository", r.Owner+"/"+r.Name),
			zap.Error(err),
		)
		return nil
	}

	return &v
}

func EncodePermissions(p *[]permission, r *Repository) string {
	data, err := json.Marshal(p)
	if err != nil {
		GetLogger().Warn(
			"Could not encode perimssions",
			zap.String("repository", r.Owner+"/"+r.Name),
			zap.Error(err),
		)
		return ""
	}
	return string(data)
}
