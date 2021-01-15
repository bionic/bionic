package twitter

import (
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"strings"
)

type Like struct {
	gorm.Model
	TweetID     int    `json:"tweetId,string" gorm:"unique"`
	FullText    string `json:"fullText"`
	ExpandedURL string `json:"expandedUrl"`
}

func (Like) TableName() string {
	return "twitter_likes"
}

func (l *Like) UnmarshalJSON(b []byte) error {
	type alias Like

	var data struct {
		Like alias `json:"like"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*l = Like(data.Like)

	return nil
}

func (p *twitter) importLikes(inputPath string) error {
	var likes []Like

	bytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	data := string(bytes)
	data = strings.TrimPrefix(data, "window.YTD.like.part0 = ")

	if err := json.Unmarshal([]byte(data), &likes); err != nil {
		return err
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(likes).
		Error
	if err != nil {
		return err
	}

	return nil
}
