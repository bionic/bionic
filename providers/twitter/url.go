package twitter

import "gorm.io/gorm"

type URL struct {
	gorm.Model
	URL      string `json:"url" gorm:"unique"`
	Expanded string `json:"expanded"`
	Display  string `json:"display"`
}

func (URL) TableName() string {
	return "twitter_urls"
}
