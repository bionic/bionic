package twitter

import (
	"encoding/json"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"path"
	"strings"
)

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
