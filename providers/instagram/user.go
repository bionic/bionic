package instagram

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
}

func (User) TableName() string {
	return tablePrefix + "users"
}

func (u User) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"username": u.Username,
	}
}
