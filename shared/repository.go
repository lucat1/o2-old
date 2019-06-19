package shared

// Repository represents a git repo in the database
type Repository struct {
	Owner       string
	Name  	    string
	Description string
	MainBranch  string `gorm:"default:'master'"`
}
