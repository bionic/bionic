package instagram

import (
	"encoding/json"
	"errors"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
)

type CommentTarget string

const (
	CommentMedia CommentTarget = "media"
)

type Comment struct {
	gorm.Model
	Target          CommentTarget `gorm:"uniqueIndex:instagram_comments_key"`
	UserID          uint          `gorm:"uniqueIndex:instagram_comments_key"`
	User            User
	Text            string           `gorm:"uniqueIndex:instagram_comments_key"`
	UserMentions    []UserMention    `gorm:"polymorphic:Parent"`
	HashtagMentions []HashtagMention `gorm:"polymorphic:Parent"`
	Timestamp       types.DateTime   `gorm:"uniqueIndex:instagram_comments_key"`
}

func (Comment) TableName() string {
	return tablePrefix + "comments"
}

func (c Comment) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"target":    c.Target,
		"user_id":   c.User.ID,
		"text":      c.Text,
		"timestamp": c.Timestamp,
	}
}

func (c *Comment) UnmarshalJSON(b []byte) error {
	var data []string

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if len(data) != 3 {
		return errors.New("incorrect comment format")
	}

	if err := c.Timestamp.UnmarshalText([]byte(data[0])); err != nil {
		return err
	}

	c.Text = data[1]
	c.User.Username = data[2]

	c.UserMentions = extractUserMentionsFromText(c.Text)
	c.HashtagMentions = extractHashtagMentionsFromText(c.Text)

	return nil
}

func (p *instagram) importComments(inputPath string) error {
	var data struct {
		MediaComments []Comment `json:"media_comments"`
	}

	bytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	for i := range data.MediaComments {
		mediaComment := &data.MediaComments[i]

		mediaComment.Target = CommentMedia

		err = p.DB().
			FirstOrCreate(&mediaComment.User, mediaComment.User.Conditions()).
			Error
		if err != nil {
			return err
		}

		err = p.DB().
			Find(mediaComment, mediaComment.Conditions()).
			Error
		if err != nil {
			return err
		}

		for j := range mediaComment.UserMentions {
			userMention := &mediaComment.UserMentions[j]

			userMention.ParentID = mediaComment.ID
			userMention.ParentType = mediaComment.TableName()

			err = p.DB().
				FirstOrCreate(&userMention.User, userMention.User.Conditions()).
				Error
			if err != nil {
				return err
			}

			err = p.DB().
				FirstOrCreate(userMention, userMention.Conditions()).
				Error
			if err != nil {
				return err
			}
		}

		for j := range mediaComment.HashtagMentions {
			hashtagMention := &mediaComment.HashtagMentions[j]

			hashtagMention.ParentID = mediaComment.ID
			hashtagMention.ParentType = mediaComment.TableName()

			err = p.DB().
				FirstOrCreate(&hashtagMention.Hashtag, hashtagMention.Hashtag.Conditions()).
				Error
			if err != nil {
				return err
			}

			err = p.DB().
				FirstOrCreate(hashtagMention, hashtagMention.Conditions()).
				Error
			if err != nil {
				return err
			}
		}

		err = p.DB().
			Clauses(clause.OnConflict{DoNothing: true}).
			Create(mediaComment).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}
