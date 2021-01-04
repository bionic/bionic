package twitter

import "gorm.io/gorm"

type Hashtag struct {
	gorm.Model
	Text string `json:"text" gorm:"unique"`
}

func (Hashtag) TableName() string {
	return "twitter_hashtags"
}
