package twitter

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID         int    `json:"id,string"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

func (User) TableName() string {
	return tablePrefix + "users"
}
