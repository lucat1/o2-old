package shared

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// GrantAccess grants the access to a certain user on a certain
// resource of a certain scope
func GrantAccess(user *User, scopes []string, owner string, repo string) error {
	GetLogger().Info(
		"Granting access to user for resource",
		zap.String("for", user.Username),
		zap.Strings("resource", []string{owner, repo}),
		zap.Strings("scopes", scopes),
	)

	var dbRepo Repository
	err := GetDatabase().Find(&dbRepo, "owner = ? AND name = ?", owner, repo).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			GetLogger().Warn("Unkown error while checking user's access", zap.Error(err))
		}

		return err
	}

	permissions := ParsePermissions(&dbRepo)

	for _, scope := range scopes {
		pex := permission{For: user.Username, Key: scope}
		v := append(*permissions, pex)
		permissions = &v
	}

	err = GetDatabase().Model(&dbRepo).Update("Granted", EncodePermissions(permissions, &dbRepo)).Error
	return err
}
