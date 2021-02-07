package telegram

import (
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"os"
	"strings"
)

type Chat struct {
	gorm.Model
	ID       int `json:"id" gorm:"unique"`
	Messages []Message
	Name     string `json:"name"`
	Type     string `json:"type"`
}

func (Chat) TableName() string {
	return tablePrefix + "chats"
}

type Message struct {
	gorm.Model
	ChatID int
	Chat   Chat

	Action                    string `json:"action"`
	Actor                     string `json:"actor"`
	ActorID                   int    `json:"actor_id"`
	Address                   string `json:"address"`
	Author                    string `json:"author"`
	ContactVcard              string `json:"contact_vcard"`
	Date                      string `json:"date"`
	DiscardReason             string `json:"discard_reason"`
	Duration                  int    `json:"duration"`
	DurationSeconds           int    `json:"duration_seconds"`
	Edited                    string `json:"edited"`
	File                      string `json:"file"`
	ForwardedFrom             string `json:"forwarded_from"`
	From                      string `json:"from"`
	FromID                    int    `json:"from_id"`
	GameDescription           string `json:"game_description"`
	GameLink                  string `json:"game_link"`
	GameMessageID             int    `json:"game_message_id"`
	GameTitle                 string `json:"game_title"`
	Height                    int    `json:"height"`
	ID                        int    `json:"id"`
	Inviter                   string `json:"inviter"`
	LiveLocationPeriodSeconds int    `json:"live_location_period_seconds"`
	MediaType                 string `json:"media_type"`
	MessageID                 int    `json:"message_id" gorm:"unique"`
	MimeType                  string `json:"mime_type"`
	Performer                 string `json:"performer"`
	Photo                     string `json:"photo"`
	PlaceName                 string `json:"place_name"`
	ReplyToMessageID          int    `json:"reply_to_message_id"`
	SavedFrom                 string `json:"saved_from"`
	Score                     int    `json:"score"`
	SelfDestructPeriodSeconds int    `json:"self_destruct_period_seconds"`
	StickerEmoji              string `json:"sticker_emoji"`
	Thumbnail                 string `json:"thumbnail"`
	Title                     string `json:"title"`
	Type                      string `json:"type"`
	ViaBot                    string `json:"via_bot"`
	Width                     int    `json:"width"`

	Text            string // todo: links, mentions
	TextAttachments []TextAttachment

	ContactInformationFirstName   string
	ContactInformationLastName    string
	ContactInformationPhoneNumber string

	LocationInformationLatitude  float64
	LocationInformationLongitude float64

	PollClosed      bool
	PollQuestion    string
	PollTotalVoters int
	PollAnswers     []PollAnswer
	Members         []Member
}

func (Message) TableName() string {
	return tablePrefix + "messages"
}

func (m *Message) UnmarshalJSON(b []byte) error {
	type alias Message

	var data struct {
		alias
		Text               TextWithAttachments `json:"text"`
		ContactInformation struct {
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			PhoneNumber string `json:"phone_number"`
		} `json:"contact_information"`
		LocationInformation struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location_information"`
		Poll struct {
			Answers []struct {
				Chosen bool   `json:"chosen"`
				Text   string `json:"text"`
				Voters int    `json:"voters"`
			} `json:"answers"`
			Closed      bool   `json:"closed"`
			Question    string `json:"question"`
			TotalVoters int    `json:"total_voters"`
		} `json:"poll"`
		Members []string `json:"members"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*m = Message(data.alias)
	m.Text = data.Text.Text
	m.TextAttachments = data.Text.Attachments
	m.ContactInformationFirstName = data.ContactInformation.FirstName
	m.ContactInformationLastName = data.ContactInformation.LastName
	m.ContactInformationPhoneNumber = data.ContactInformation.PhoneNumber
	m.LocationInformationLatitude = data.LocationInformation.Latitude
	m.LocationInformationLongitude = data.LocationInformation.Longitude
	m.PollClosed = data.Poll.Closed
	m.PollQuestion = data.Poll.Question
	m.PollTotalVoters = data.Poll.TotalVoters
	for _, answer := range data.Poll.Answers {
		m.PollAnswers = append(
			m.PollAnswers,
			PollAnswer{Chosen: answer.Chosen, Text: answer.Text, Voters: answer.Voters})

	}
	for _, member := range data.Members {
		m.Members = append(m.Members, Member{Name: member})
	}

	return nil
}

type TextAttachment struct {
	gorm.Model
	MessageID int
	Message   Message
	Type      string `json:"type"`
	Text      string `json:"text"`
	Href      string `json:"href"`
}

func (TextAttachment) TableName() string {
	return tablePrefix + "text_attachments"
}

type PollAnswer struct {
	gorm.Model
	MessageID int
	Message   Message
	Chosen    bool
	Text      string
	Voters    int
}

func (PollAnswer) TableName() string {
	return tablePrefix + "poll_answers"
}

type Member struct {
	gorm.Model
	MessageID int
	Message   Message
	Name      string `json:"name"`
}

func (Member) TableName() string {
	return tablePrefix + "members"
}

func (p *telegram) importChats(inputPath string) error {
	rc, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var data struct {
		Chats struct {
			About string `json:"about"`
			List  []Chat `json:"list"`
		} `json:"chats"`
	}

	bytes, err := io.ReadAll(rc)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	err = p.DB().
		Session(&gorm.Session{CreateBatchSize: 100}).
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(&data.Chats.List).
		Error
	if err != nil {
		return err
	}

	return nil
}

type TextWithAttachments struct {
	Text        string
	Attachments []TextAttachment
}

func (t *TextWithAttachments) UnmarshalJSON(b []byte) error {
	// The "text" field could be one of two types:
	// 1) 'text': ['btc\n', {'type': 'link', 'text': 'coin.space'}, '\nabc123']
	// 2) 'text': 'one two three'
	//
	// For the first type, we save all attachments as related objects and build text as concatenation of text parts
	// (like "btc\n") and 'text' params of attachments (like "coin.space").
	//
	// Sometimes attachment object has 'href' parameter. In this case, we inject the link into text like this:
	// "google (https://google.com/)"

	var text string
	err := json.Unmarshal(b, &text)
	if err == nil {
		t.Text = text
		return nil
	}

	var partsOrAttachments []TextPartOrTextAttachment
	err = json.Unmarshal(b, &partsOrAttachments)
	if err != nil {
		return err
	}

	var textBuilder strings.Builder
	var attachments []TextAttachment
	for _, partOrAttachment := range partsOrAttachments {
		if partOrAttachment.isAttachment {
			textBuilder.WriteString(partOrAttachment.Attachment.Text)
			if partOrAttachment.Attachment.Href != "" {
				textBuilder.WriteString(" (" + partOrAttachment.Attachment.Href + ")")
			}
			attachments = append(attachments, partOrAttachment.Attachment)
		} else {
			textBuilder.WriteString(partOrAttachment.Part)
		}
	}

	t.Text = textBuilder.String()
	t.Attachments = attachments

	return nil
}

type TextPartOrTextAttachment struct {
	isAttachment bool
	Part         string
	Attachment   TextAttachment
}

func (t *TextPartOrTextAttachment) UnmarshalJSON(b []byte) error {
	var part string
	err := json.Unmarshal(b, &part)
	if err == nil {
		t.Part = part
		return nil
	}

	var attachment TextAttachment
	err = json.Unmarshal(b, &attachment)
	if err != nil {
		return err
	}

	t.isAttachment = true
	t.Attachment = attachment

	return nil
}
