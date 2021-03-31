package twitter

import (
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Like struct {
	gorm.Model
	TweetID     int `json:"tweetId,string" gorm:"unique"`
	Tweet       Tweet
	FullText    string `json:"fullText"`
	ExpandedURL string `json:"expandedUrl"`
}

func (Like) TableName() string {
	return tablePrefix + "likes"
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

	if err := readJSON(
		inputPath,
		"window.YTD.like.part0 = ",
		&likes,
	); err != nil {
		return err
	}

	for i, like := range likes {
		err := p.DB().
			FirstOrCreate(&likes[i].Tweet, map[string]interface{}{
				"id": like.TweetID,
			}).
			Error
		if err != nil {
			return err
		}
	}

	err := p.DB().
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
