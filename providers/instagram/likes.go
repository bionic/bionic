package instagram

import (
	"encoding/json"
	"errors"
	"github.com/bionic-dev/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
)

type LikeTarget string

const (
	LikeMedia   LikeTarget = "media"
	LikeComment LikeTarget = "comment"
)

type Like struct {
	gorm.Model
	Target    LikeTarget `gorm:"uniqueIndex:instagram_likes_key"`
	UserID    uint       `gorm:"uniqueIndex:instagram_likes_key"`
	User      User
	Timestamp types.DateTime `gorm:"uniqueIndex:instagram_likes_key"`
}

func (Like) TableName() string {
	return tablePrefix + "likes"
}

func (l *Like) UnmarshalJSON(b []byte) error {
	var data []string

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if len(data) != 2 {
		return errors.New("incorrect like format")
	}

	if err := l.Timestamp.UnmarshalText([]byte(data[0])); err != nil {
		return err
	}

	l.User.Username = data[1]

	return nil
}

func (p *instagram) importLikes(inputPath string) error {
	var data struct {
		MediaLikes   []Like `json:"media_likes"`
		CommentLikes []Like `json:"comment_likes"`
	}

	bytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	for i := range data.MediaLikes {
		mediaLike := &data.MediaLikes[i]

		mediaLike.Target = LikeMedia

		err = p.DB().
			FirstOrCreate(&mediaLike.User, mediaLike.User.Conditions()).
			Error
		if err != nil {
			return err
		}
	}

	for i := range data.CommentLikes {
		commentLike := &data.CommentLikes[i]

		commentLike.Target = LikeComment

		err = p.DB().
			FirstOrCreate(&commentLike.User, commentLike.User.Conditions()).
			Error
		if err != nil {
			return err
		}
	}

	err = p.DB().
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(data.MediaLikes, 100).
		Error
	if err != nil {
		return err
	}

	err = p.DB().
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(data.CommentLikes, 100).
		Error
	if err != nil {
		return err
	}

	return nil
}
