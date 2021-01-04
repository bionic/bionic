package twitter

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	TweetID     int    `json:"tweetId,string" gorm:"unique"`
	FullText    string `json:"fullText"`
	ExpandedURL string `json:"expandedUrl"`
}

func (Like) TableName() string {
	return "twitter_likes"
}
