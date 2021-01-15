package twitter

import (
	"encoding/json"
	"github.com/shekhirin/bionic-cli/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"strconv"
	"strings"
)

type Tweet struct {
	gorm.Model
	ID                 int `json:"id,string"`
	AuthorID           *int
	Author             *User
	Retweeted          bool          `json:"retweeted"`
	Source             string        `json:"source"`
	Entities           TweetEntities `json:"entities"`
	DisplayTextFromIdx *int
	DisplayTextToIdx   *int
	FavoriteCount      int            `json:"favorite_count,string"`
	Truncated          bool           `json:"truncated"`
	RetweetCount       int            `json:"retweet_count,string"`
	PossiblySensitive  bool           `json:"possibly_sensitive"`
	Created            types.DateTime `json:"created_at"`
	Favorited          bool           `json:"favorited"`
	FullText           string         `json:"full_text"`
	Lang               string         `json:"lang"`
	InReplyToUserID    *int
	InReplyToUser      *User
	InReplyToStatusID  *int
	InReplyToStatus    *Tweet
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
	FromIdx         *int
	ToIdx           *int
}

func (TweetHashtag) TableName() string {
	return "twitter_tweet_hashtags"
}

type TweetMedia struct {
	gorm.Model
	TweetEntitiesID int
	ID              int    `json:"id,string"`
	ExpandedURL     string `json:"expanded_url"`
	FromIdx         *int
	ToIdx           *int
	URL             string `json:"url"`
	MediaURL        string `json:"media_url"`
	MediaURLHTTPS   string `json:"media_url_https"`
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
	FromIdx         *int
	ToIdx           *int
}

func (TweetUserMention) TableName() string {
	return "twitter_tweet_user_mentions"
}

type TweetURL struct {
	gorm.Model
	TweetEntitiesID int
	URLID           string
	URL             URL
	FromIdx         *int
	ToIdx           *int
}

func (TweetURL) TableName() string {
	return "twitter_tweet_urls"
}

func (p *twitter) importTweets(inputPath string) error {
	// TODO: use json.Unmarshaler as in personalization.go
	var fileData []struct {
		Tweet struct {
			Tweet
			DisplayTextRange []string `json:"display_text_range"`
			Entities         struct {
				TweetEntities
				Hashtags []struct {
					TweetHashtag
					Hashtag
					Indices []string `json:"indices"`
				} `json:"hashtags"`
				Media []struct {
					TweetMedia
					Indices []string `json:"indices"`
				} `json:"media"`
				UserMentions []struct {
					TweetUserMention
					User
					Indices []string `json:"indices"`
				} `json:"user_mentions"`
				URLs []struct {
					TweetURL
					URL
					Indices  []string `json:"indices"`
					Expanded string   `json:"expanded_url"`
					Display  string   `json:"display_url"`
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

		tweet.DisplayTextFromIdx, tweet.DisplayTextToIdx = rangeToIndices(entry.Tweet.DisplayTextRange)

		tweet.Entities = entry.Tweet.Entities.TweetEntities

		for _, hashtag := range entry.Tweet.Entities.Hashtags {
			tweetHashtag := hashtag.TweetHashtag
			tweetHashtag.Hashtag = hashtag.Hashtag
			tweetHashtag.FromIdx, tweetHashtag.ToIdx = rangeToIndices(hashtag.Indices)

			tweet.Entities.Hashtags = append(tweet.Entities.Hashtags, tweetHashtag)
		}

		for _, media := range entry.Tweet.Entities.Media {
			tweetMedia := media.TweetMedia
			tweetMedia.FromIdx, tweetMedia.ToIdx = rangeToIndices(media.Indices)

			tweet.Entities.Media = append(tweet.Entities.Media, tweetMedia)
		}

		for _, userMention := range entry.Tweet.Entities.UserMentions {
			tweetUserMention := userMention.TweetUserMention
			tweetUserMention.User = userMention.User
			tweetUserMention.FromIdx, tweetUserMention.ToIdx = rangeToIndices(userMention.Indices)

			tweet.Entities.UserMentions = append(tweet.Entities.UserMentions, tweetUserMention)
		}

		for _, url := range entry.Tweet.Entities.URLs {
			tweetURL := url.TweetURL
			tweetURL.URL = url.URL
			tweetURL.URL.Expanded = url.Expanded
			tweetURL.URL.Display = url.Display
			tweetURL.FromIdx, tweetURL.ToIdx = rangeToIndices(url.Indices)

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

func rangeToIndices(indicesRange []string) (*int, *int) {
	if len(indicesRange) != 2 {
		return nil, nil
	}

	from, to := indicesRange[0], indicesRange[1]

	fromInt, err := strconv.Atoi(from)
	if err != nil {
		return nil, nil
	}

	toInt, err := strconv.Atoi(to)
	if err != nil {
		return nil, nil
	}

	return &fromInt, &toInt
}
