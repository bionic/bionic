package twitter

import (
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

type Conversation struct {
	gorm.Model
	ID             string `json:"conversationId"`
	DirectMessages []DirectMessage
}

func (c Conversation) TableName() string {
	return "twitter_conversations"
}

type DirectMessage struct {
	gorm.Model
	ConversationID string
	ID             string     `json:"id"`
	RecipientID    string     `json:"recipientId"`
	Reactions      []Reaction `json:"reactions"`
	URLs           []URL      `json:"urls"`
	Text           string     `json:"text"`
	MediaURLs      []string   `json:"mediaUrls" gorm:"type:text"`
	SenderID       string     `json:"senderId"`
	Created        time.Time  `json:"createdAt"`
}

func (dm DirectMessage) TableName() string {
	return "twitter_direct_messages"
}

type Reaction struct {
	gorm.Model
	DirectMessageID string
	SenderID        string    `json:"senderId"`
	ReactionKey     string    `json:"reactionKey"`
	EventID         string    `json:"eventId"`
	Created         time.Time `json:"createdAt"`
}

func (r Reaction) TableName() string {
	return "twitter_reactions"
}

type URL struct {
	gorm.Model
	DirectMessageID string
	URL             string `json:"url"`
	Expanded        string `json:"expanded"`
	Display         string `json:"display"`
}

func (u URL) TableName() string {
	return "twitter_urls"
}

func (p *twitter) processDirectMessages(inputPath string) error {
	var fileData []struct {
		DmConversation struct {
			Conversation
			Messages []struct {
				MessageCreate DirectMessage `json:"messageCreate"`
			} `json:"messages"`
		} `json:"dmConversation"`
	}

	bytes, err := ioutil.ReadFile(path.Join(inputPath, "data", "direct-messages.js"))
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

	err = p.db.
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
