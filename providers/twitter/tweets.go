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
	return tablePrefix + "tweets"
}

func (t *Tweet) UnmarshalJSON(b []byte) error {
	type alias Tweet

	var data struct {
		Tweet struct {
			alias
			DisplayTextRange    []string `json:"display_text_range"`
			InReplyToStatusID   *int     `json:"in_reply_to_status_id,string"`
			InReplyToUserID     *int     `json:"in_reply_to_user_id,string"`
			InReplyToScreenName *string  `json:"in_reply_to_screen_name"`
		} `json:"tweet"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	tweet := data.Tweet

	*t = Tweet(tweet.alias)

	t.DisplayTextFromIdx, t.DisplayTextToIdx = rangeToIndices(tweet.DisplayTextRange)

	if tweet.InReplyToStatusID != nil {
		t.InReplyToStatus = &Tweet{
			ID: *tweet.InReplyToStatusID,
		}
	}

	if tweet.InReplyToUserID != nil && tweet.InReplyToScreenName != nil {
		t.InReplyToUser = &User{
			ID:         *tweet.InReplyToUserID,
			ScreenName: *tweet.InReplyToScreenName,
		}

		if t.InReplyToStatus != nil {
			t.InReplyToStatus.Author = t.InReplyToUser
		}
	}

	return nil
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
	return tablePrefix + "tweet_entities"
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
	return tablePrefix + "tweet_hashtags"
}

func (th *TweetHashtag) UnmarshalJSON(b []byte) error {
	type alias TweetHashtag

	var data struct {
		alias
		Hashtag
		Indices []string `json:"indices"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*th = TweetHashtag(data.alias)

	th.Hashtag = data.Hashtag
	th.FromIdx, th.ToIdx = rangeToIndices(data.Indices)

	return nil
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
	return tablePrefix + "tweet_media"
}

func (tm *TweetMedia) UnmarshalJSON(b []byte) error {
	type alias TweetMedia

	var data struct {
		alias
		Indices []string `json:"indices"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*tm = TweetMedia(data.alias)

	tm.FromIdx, tm.ToIdx = rangeToIndices(data.Indices)

	return nil
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
	return tablePrefix + "tweet_user_mentions"
}

func (tum *TweetUserMention) UnmarshalJSON(b []byte) error {
	type alias TweetUserMention

	var data struct {
		alias
		User
		Indices []string `json:"indices"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*tum = TweetUserMention(data.alias)

	tum.User = data.User
	tum.FromIdx, tum.ToIdx = rangeToIndices(data.Indices)

	return nil
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
	return tablePrefix + "tweet_urls"
}

func (tu *TweetURL) UnmarshalJSON(b []byte) error {
	type alias TweetURL

	var data struct {
		alias
		URL
		Indices  []string `json:"indices"`
		Expanded string   `json:"expanded_url"`
		Display  string   `json:"display_url"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*tu = TweetURL(data.alias)

	tu.URL = data.URL
	tu.URL.Expanded = data.Expanded
	tu.URL.Display = data.Display
	tu.FromIdx, tu.ToIdx = rangeToIndices(data.Indices)

	return nil
}

func (p *twitter) importTweets(inputPath string) error {
	var tweets []Tweet

	bytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	data := string(bytes)
	data = strings.TrimPrefix(data, "window.YTD.tweet.part0 = ")

	if err := json.Unmarshal([]byte(data), &tweets); err != nil {
		return err
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
