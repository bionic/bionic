package twitter

import (
	"encoding/json"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"strings"
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

func (p *twitter) importTweets(inputPath string) error {
	var fileData []struct {
		Tweet struct {
			Tweet
			Entities struct {
				TweetEntities
				Hashtags []struct {
					TweetHashtag
					Hashtag
				} `json:"hashtags"`
				UserMentions []struct {
					TweetUserMention
					User
				} `json:"user_mentions"`
				URLs []struct {
					TweetURL
					URL
					Expanded string `json:"expanded_url"`
					Display  string `json:"display_url"`
				} `json:"urls"`
			} `json:"entities"`
			InReplyToStatusID   *int    `json:"in_reply_to_status_id,string"`
			InReplyToUserID     *int    `json:"in_reply_to_user_id,string"`
			InReplyToScreenName *string `json:"in_reply_to_screen_name"`
		} `json:"tweet"`
	}

	bytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	data := string(bytes)
	data = strings.TrimPrefix(data, "window.YTD.tweet.part0 = ")

	if err := json.Unmarshal([]byte(data), &fileData); err != nil {
		return err
	}

	var tweets []Tweet

	for _, entry := range fileData {
		tweet := entry.Tweet.Tweet

		tweet.Entities = entry.Tweet.Entities.TweetEntities

		for _, hashtag := range entry.Tweet.Entities.Hashtags {
			tweetHashtag := hashtag.TweetHashtag
			tweetHashtag.Hashtag = hashtag.Hashtag
			tweet.Entities.Hashtags = append(tweet.Entities.Hashtags, tweetHashtag)
		}

		for _, userMention := range entry.Tweet.Entities.UserMentions {
			tweetUserMention := userMention.TweetUserMention
			tweetUserMention.User = userMention.User
			tweet.Entities.UserMentions = append(tweet.Entities.UserMentions, tweetUserMention)
		}

		for _, url := range entry.Tweet.Entities.URLs {
			tweetURL := url.TweetURL
			tweetURL.URL = url.URL
			tweetURL.URL.Expanded = url.Expanded
			tweetURL.URL.Display = url.Display
			tweet.Entities.URLs = append(tweet.Entities.URLs, tweetURL)
		}

		if entry.Tweet.InReplyToStatusID != nil {
			tweet.InReplyToStatus = &Tweet{
				ID: *entry.Tweet.InReplyToStatusID,
			}
		}

		if entry.Tweet.InReplyToUserID != nil && entry.Tweet.InReplyToScreenName != nil {
			tweet.InReplyToUser = &User{
				ID:         *entry.Tweet.InReplyToUserID,
				ScreenName: *entry.Tweet.InReplyToScreenName,
			}

			if tweet.InReplyToStatus != nil {
				tweet.InReplyToStatus.Author = tweet.InReplyToUser
			}
		}

		tweets = append(tweets, tweet)
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(tweets).
		Error
	if err != nil {
		return err
	}

	return nil
}
