package twitter

import (
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"strings"
	"time"
)

type Conversation struct {
	gorm.Model
	ID             string          `json:"conversationId"`
	DirectMessages []DirectMessage `json:"messages"`
}

func (Conversation) TableName() string {
	return tablePrefix + "conversations"
}

func (c *Conversation) UnmarshalJSON(b []byte) error {
	type alias Conversation

	var data struct {
		DmConversation alias `json:"dmConversation"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*c = Conversation(data.DmConversation)

	return nil
}

type DirectMessage struct {
	gorm.Model
	ConversationID string
	ID             int                     `json:"id,string"`
	RecipientID    int                     `json:"recipientId,string"`
	Reactions      []DirectMessageReaction `json:"reactions"`
	URLs           []URL                   `json:"urls" gorm:"many2many:twitter_direct_message_urls"`
	Text           string                  `json:"text"`
	MediaURLs      []string                `json:"mediaUrls" gorm:"type:text"`
	SenderID       int                     `json:"senderId,string"`
	Created        time.Time               `json:"createdAt"`
}

func (DirectMessage) TableName() string {
	return tablePrefix + "direct_messages"
}

func (dm *DirectMessage) UnmarshalJSON(b []byte) error {
	type alias DirectMessage

	var data struct {
		MessageCreate alias `json:"messageCreate"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*dm = DirectMessage(data.MessageCreate)

	return nil
}

type DirectMessageReaction struct {
	gorm.Model
	DirectMessageID int
	SenderID        string    `json:"senderId"`
	ReactionKey     string    `json:"reactionKey"`
	EventID         string    `json:"eventId"`
	Created         time.Time `json:"createdAt"`
}

func (DirectMessageReaction) TableName() string {
	return tablePrefix + "direct_message_reactions"
}

func (p *twitter) importDirectMessages(inputPath string) error {
	var conversations []Conversation

	bytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	data := string(bytes)
	data = strings.TrimPrefix(data, "window.YTD.direct_messages.part0 = ")

	if err := json.Unmarshal([]byte(data), &conversations); err != nil {
		return err
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(conversations).
		Error
	if err != nil {
		return err
	}

	return nil
}
