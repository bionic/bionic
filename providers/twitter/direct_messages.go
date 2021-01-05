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
	ID             string `json:"conversationId"`
	DirectMessages []DirectMessage
}

func (Conversation) TableName() string {
	return "twitter_conversations"
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
	return "twitter_direct_messages"
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
	return "twitter_direct_message_reactions"
}

func (p *twitter) importDirectMessages(inputPath string) error {
	var fileData []struct {
		DmConversation struct {
			Conversation
			Messages []struct {
				MessageCreate DirectMessage `json:"messageCreate"`
			} `json:"messages"`
		} `json:"dmConversation"`
	}

	bytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	data := string(bytes)
	data = strings.TrimPrefix(data, "window.YTD.direct_messages.part0 = ")

	if err := json.Unmarshal([]byte(data), &fileData); err != nil {
		return err
	}

	var conversations []Conversation

	for _, entry := range fileData {
		conversation := entry.DmConversation.Conversation
		for _, message := range entry.DmConversation.Messages {
			conversation.DirectMessages = append(conversation.DirectMessages, message.MessageCreate)
		}

		conversations = append(conversations, conversation)
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
