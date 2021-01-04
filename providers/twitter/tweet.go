package twitter

import (
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
)

type Tweet struct {
	gorm.Model
	ID                int `json:"id,string"`
	AuthorID          *int
	Author            *User
	Retweeted         bool                 `json:"retweeted"`
	Source            string               `json:"source"`
	Entities          TweetEntities        `json:"entities"`
	DisplayTextRange  types.IntStringSlice `json:"display_text_range"`
	FavoriteCount     int                  `json:"favorite_count,string"`
	Truncated         bool                 `json:"truncated"`
	RetweetCount      int                  `json:"retweet_count,string"`
	PossiblySensitive bool                 `json:"possibly_sensitive"`
	Created           types.DateTime       `json:"created_at"`
	Favorited         bool                 `json:"favorited"`
	FullText          string               `json:"full_text"`
	Lang              string               `json:"lang"`
	InReplyToUserID   *int
	InReplyToUser     *User
	InReplyToStatusID *int
	InReplyToStatus   *Tweet
}

func (Tweet) TableName() string {
	return "twitter_tweets"
}

type TweetEntities struct {
	gorm.Model
	TweetID  int
	Hashtags []TweetHashtag `json:"hashtags"`
	Media    []TweetMedia   `json:"media"`
	//Symbols      []Symbol       `json:"symbols"`
	UserMentions []TweetUserMention `json:"user_mentions"`
	URLs         []TweetURL         `json:"urls"`
}

func (TweetEntities) TableName() string {
	return "twitter_tweet_entities"
}

type TweetHashtag struct {
	gorm.Model
	TweetEntitiesID int
	HashtagID       int
	Hashtag         Hashtag
	Indices         types.IntStringSlice `json:"indices"`
}

func (TweetHashtag) TableName() string {
	return "twitter_tweet_hashtags"
}

type TweetMedia struct {
	gorm.Model
	TweetEntitiesID int
	ID              int                  `json:"id,string"`
	ExpandedURL     string               `json:"expanded_url"`
	Indices         types.IntStringSlice `json:"indices"`
	URL             string               `json:"url"`
	MediaURL        string               `json:"media_url"`
	MediaURLHTTPS   string               `json:"media_url_https"`
	//Sizes           struct {
	//	Thumb struct {
	//		W      string `json:"w"`
	//		H      string `json:"h"`
	//		Resize string `json:"resize"`
	//	} `json:"thumb"`
	//	Small struct {
	//		W      string `json:"w"`
	//		H      string `json:"h"`
	//		Resize string `json:"resize"`
	//	} `json:"small"`
	//	Large struct {
	//		W      string `json:"w"`
	//		H      string `json:"h"`
	//		Resize string `json:"resize"`
	//	} `json:"large"`
	//	Medium struct {
	//		W      string `json:"w"`
	//		H      string `json:"h"`
	//		Resize string `json:"resize"`
	//	} `json:"medium"`
	//} `json:"sizes"`
	Type       string `json:"type"`
	DisplayURL string `json:"display_url"`
}

func (TweetMedia) TableName() string {
	return "twitter_tweet_media"
}

type TweetUserMention struct {
	gorm.Model
	TweetEntitiesID int
	UserID          int
	User            User
	Indices         types.IntStringSlice `json:"indices"`
}

func (TweetUserMention) TableName() string {
	return "twitter_tweet_user_mentions"
}

type TweetURL struct {
	gorm.Model
	TweetEntitiesID int
	URLID           string
	URL             URL
	Indices         types.IntStringSlice `json:"indices"`
}

func (TweetURL) TableName() string {
	return "twitter_tweet_urls"
}
