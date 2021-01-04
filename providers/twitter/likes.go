package twitter

import (
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"path"
	"strings"
)

type Like struct {
	gorm.Model
	TweetID     string `json:"tweetId" gorm:"unique"`
	FullText    string `json:"fullText"`
	ExpandedURL string `json:"expandedUrl"`
}

func (l Like) TableName() string {
	return "twitter_likes"
}

func (p *twitter) processLikes(inputPath string) error {
	var fileData []struct {
		Like Like `json:"like"`
	}

	bytes, err := ioutil.ReadFile(path.Join(inputPath, "data", "like.js"))
	if err != nil {
		return err
	}

	data := string(bytes)
	data = strings.TrimPrefix(data, "window.YTD.like.part0 = ")

	if err := json.Unmarshal([]byte(data), &fileData); err != nil {
		return err
	}

	var likes []Like

	for _, entry := range fileData {
		likes = append(likes, entry.Like)
	}

	err = p.db.
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
