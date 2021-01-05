package twitter

import (
	"encoding/json"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"strings"
)

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
