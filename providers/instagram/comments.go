package instagram

import (
	"encoding/json"
	"errors"
	"github.com/bionic-dev/bionic/types"
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
	UserID          int           `gorm:"uniqueIndex:instagram_comments_key"`
	User            User
	Text            string `gorm:"uniqueIndex:instagram_comments_key"`
	UserMentions    []CommentUserMention
	HashtagMentions []CommentHashtagMention
	Timestamp       types.DateTime `gorm:"uniqueIndex:instagram_comments_key"`
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

	for _, userMention := range extractUserMentionsFromText(c.Text) {
		c.UserMentions = append(c.UserMentions, CommentUserMention{
			User: User{
				Username: userMention.Username,
			},
			FromIdx: userMention.FromIdx,
			ToIdx:   userMention.ToIdx,
		})
	}

	for _, hashtagMention := range extractHashtagMentionsFromText(c.Text) {
		c.HashtagMentions = append(c.HashtagMentions, CommentHashtagMention{
			Hashtag: Hashtag{
				Text: hashtagMention.Hashtag,
			},
			FromIdx: hashtagMention.FromIdx,
			ToIdx:   hashtagMention.ToIdx,
		})
	}

	return nil
}

type CommentUserMention struct {
	gorm.Model
	CommentID uint `gorm:"uniqueIndex:instagram_comment_user_mentions_key"`
	UserID    int  `gorm:"uniqueIndex:instagram_comment_user_mentions_key"`
	User      User
	FromIdx   int `gorm:"uniqueIndex:instagram_comment_user_mentions_key"`
	ToIdx     int `gorm:"uniqueIndex:instagram_comment_user_mentions_key"`
}

func (CommentUserMention) TableName() string {
	return tablePrefix + "comment_user_mentions"
}

func (cum CommentUserMention) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"comment_id": cum.CommentID,
		"user_id":    cum.User.ID,
		"from_idx":   cum.FromIdx,
		"to_idx":     cum.ToIdx,
	}
}

type CommentHashtagMention struct {
	gorm.Model
	CommentID uint `gorm:"uniqueIndex:instagram_comment_hashtag_mentions_key"`
	HashtagID int  `gorm:"uniqueIndex:instagram_comment_hashtag_mentions_key"`
	Hashtag   Hashtag
	FromIdx   int `gorm:"uniqueIndex:instagram_comment_hashtag_mentions_key"`
	ToIdx     int `gorm:"uniqueIndex:instagram_comment_hashtag_mentions_key"`
}

func (CommentHashtagMention) TableName() string {
	return tablePrefix + "comment_hashtag_mentions"
}

func (chm CommentHashtagMention) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"comment_id": chm.CommentID,
		"hashtag_id": chm.Hashtag.ID,
		"from_idx":   chm.FromIdx,
		"to_idx":     chm.ToIdx,
	}
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

			userMention.CommentID = mediaComment.ID

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

			hashtagMention.CommentID = mediaComment.ID

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
