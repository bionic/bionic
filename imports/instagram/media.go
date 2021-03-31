package instagram

import (
	"encoding/json"
	"github.com/bionic-dev/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
)

type MediaType string

const (
	MediaStory MediaType = "story"
	MediaVideo MediaType = "video"
	MediaPhoto MediaType = "photo"
)

type MediaItem struct {
	gorm.Model
	Type            MediaType             `gorm:"uniqueIndex:instagram_media_key"`
	Caption         *string               `json:"caption"`
	TakenAt         types.DateTime        `json:"taken_at" gorm:"uniqueIndex:instagram_media_key"`
	Location        *string               `json:"location"`
	Path            string                `json:"path"`
	UserMentions    []MediaUserMention    `gorm:"foreignKey:MediaID"`
	HashtagMentions []MediaHashtagMention `gorm:"foreignKey:MediaID"`
}

func (MediaItem) TableName() string {
	return tablePrefix + "media"
}

func (mi MediaItem) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"type":     mi.Type,
		"taken_at": mi.TakenAt,
	}
}

func (mi *MediaItem) UnmarshalJSON(b []byte) error {
	type alias MediaItem
	var data alias

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*mi = MediaItem(data)

	if mi.Caption != nil && *mi.Caption == "" {
		mi.Caption = nil
	}

	return nil
}

type MediaUserMention struct {
	gorm.Model
	MediaID uint `gorm:"uniqueIndex:instagram_media_user_mentions_key"`
	UserID  int  `gorm:"uniqueIndex:instagram_media_user_mentions_key"`
	User    User
	FromIdx int `gorm:"uniqueIndex:instagram_media_user_mentions_key"`
	ToIdx   int `gorm:"uniqueIndex:instagram_media_user_mentions_key"`
}

func (MediaUserMention) TableName() string {
	return tablePrefix + "media_user_mentions"
}

func (mum MediaUserMention) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"media_id": mum.MediaID,
		"user_id":  mum.User.ID,
		"from_idx": mum.FromIdx,
		"to_idx":   mum.ToIdx,
	}
}

type MediaHashtagMention struct {
	gorm.Model
	MediaID   uint `gorm:"uniqueIndex:instagram_media_hashtag_mentions_key"`
	HashtagID int  `gorm:"uniqueIndex:instagram_media_hashtag_mentions_key"`
	Hashtag   Hashtag
	FromIdx   int `gorm:"uniqueIndex:instagram_media_hashtag_mentions_key"`
	ToIdx     int `gorm:"uniqueIndex:instagram_media_hashtag_mentions_key"`
}

func (MediaHashtagMention) TableName() string {
	return tablePrefix + "media_hashtag_mentions"
}

func (mhm MediaHashtagMention) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"media_id":   mhm.MediaID,
		"hashtag_id": mhm.Hashtag.ID,
		"from_idx":   mhm.FromIdx,
		"to_idx":     mhm.ToIdx,
	}
}

type ProfilePhoto struct {
	gorm.Model
	TakenAt         types.DateTime `json:"taken_at" gorm:"unique"`
	IsActiveProfile bool           `json:"is_active_profile"`
	Path            string         `json:"path"`
}

func (ProfilePhoto) TableName() string {
	return tablePrefix + "profile_photos"
}

func (p *instagram) importMedia(inputPath string) error {
	var data struct {
		Stories []MediaItem    `json:"stories"`
		Videos  []MediaItem    `json:"videos"`
		Profile []ProfilePhoto `json:"profile"`
		Photos  []MediaItem    `json:"photos"`
	}

	bytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	if len(data.Profile) > 0 {
		err = p.DB().
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "taken_at"}},
				DoUpdates: clause.AssignmentColumns([]string{"is_active_profile"}),
			}).
			Create(data.Profile).
			Error
		if err != nil {
			return err
		}
	}

	if err := p.insertMedia(data.Stories, MediaStory); err != nil {
		return err
	}

	if err := p.insertMedia(data.Videos, MediaVideo); err != nil {
		return err
	}

	if err := p.insertMedia(data.Photos, MediaPhoto); err != nil {
		return err
	}

	return nil
}

func (p *instagram) insertMedia(media []MediaItem, mediaType MediaType) error {
	for i := range media {
		media[i].Type = mediaType

		err := p.DB().
			Find(&media[i], media[i].Conditions()).
			Error
		if err != nil {
			return err
		}

		if err := p.insertMediaUserMentions(&media[i]); err != nil {
			return err
		}

		if err := p.insertMediaHashtagMentions(&media[i]); err != nil {
			return err
		}

		err = p.DB().
			Clauses(clause.OnConflict{DoNothing: true}).
			Create(&media[i]).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *instagram) insertMediaUserMentions(media *MediaItem) error {
	if media.Caption == nil {
		return nil
	}

	for _, userMention := range extractUserMentionsFromText(*media.Caption) {
		mediaUserMention := MediaUserMention{
			User: User{
				Username: userMention.Username,
			},
			FromIdx: userMention.FromIdx,
			ToIdx:   userMention.ToIdx,
		}

		mediaUserMention.MediaID = media.ID

		err := p.DB().
			FirstOrCreate(&mediaUserMention.User, mediaUserMention.User.Conditions()).
			Error
		if err != nil {
			return err
		}

		err = p.DB().
			FirstOrCreate(&mediaUserMention, mediaUserMention.Conditions()).
			Error
		if err != nil {
			return err
		}

		media.UserMentions = append(media.UserMentions, mediaUserMention)
	}

	return nil
}

func (p *instagram) insertMediaHashtagMentions(media *MediaItem) error {
	if media.Caption == nil {
		return nil
	}

	for _, hashtagMention := range extractHashtagMentionsFromText(*media.Caption) {
		mediaHashtagMention := MediaHashtagMention{
			Hashtag: Hashtag{
				Text: hashtagMention.Hashtag,
			},
			FromIdx: hashtagMention.FromIdx,
			ToIdx:   hashtagMention.ToIdx,
		}

		mediaHashtagMention.MediaID = media.ID

		err := p.DB().
			FirstOrCreate(&mediaHashtagMention.Hashtag, mediaHashtagMention.Hashtag.Conditions()).
			Error
		if err != nil {
			return err
		}

		err = p.DB().
			FirstOrCreate(&mediaHashtagMention, mediaHashtagMention.Conditions()).
			Error
		if err != nil {
			return err
		}

		media.HashtagMentions = append(media.HashtagMentions, mediaHashtagMention)
	}

	return nil
}
