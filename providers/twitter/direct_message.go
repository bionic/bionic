package twitter

import (
	"gorm.io/gorm"
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
