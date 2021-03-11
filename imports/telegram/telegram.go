package telegram

import (
	"github.com/bionic-dev/bionic/imports/provider"
	"github.com/bionic-dev/bionic/internal/provider/database"
	"gorm.io/gorm"
	"path"
)

const name = "telegram"
const tablePrefix = "telegram_"

type telegram struct {
	database.Database
}

func New(db *gorm.DB) provider.Provider {
	return &telegram{
		Database: database.New(db),
	}
}

func (telegram) Name() string {
	return name
}

func (telegram) TablePrefix() string {
	return tablePrefix
}

func (telegram) ImportDescription() string {
	return "Desktop App (https://desktop.telegram.org/ only): Settings => Advanced => Export Telegram data"
}

func (p *telegram) Migrate() error {
	err := p.DB().AutoMigrate(
		&Chat{},
		&Message{},
		&TextAttachment{},
		&PollAnswer{},
		&Member{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *telegram) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	if !provider.IsPathDir(inputPath) {
		return nil, provider.ErrInputPathShouldBeDirectory
	}

	return []provider.ImportFn{
		provider.NewImportFn(
			"Chats",
			p.importChats,
			path.Join(inputPath, "result.json"),
		),
	}, nil
}

// TODO: Add remaining structures
//type Foo struct {
//	About string `json:"about"`
//	Contacts struct {
//		About string `json:"about"`
//		List  []struct {
//			Date        string `json:"date"`
//			FirstName   string `json:"first_name"`
//			LastName    string `json:"last_name"`
//			PhoneNumber string `json:"phone_number"`
//			UserID      int64  `json:"user_id"`
//		} `json:"list"`
//	} `json:"contacts"`
//	FrequentContacts struct {
//		About string `json:"about"`
//		List  []struct {
//			Category string  `json:"category"`
//			ID       int64   `json:"id"`
//			Name     string  `json:"name"`
//			Rating   float64 `json:"rating"`
//			Type     string  `json:"type"`
//		} `json:"list"`
//	} `json:"frequent_contacts"`
//	LeftChats struct {
//		About string `json:"about"`
//		List  []struct {
//			ID       int64 `json:"id"`
//			Messages []struct {
//				Action           string   `json:"action"`
//				Actor            string   `json:"actor"`
//				ActorID          int64    `json:"actor_id"`
//				Date             string   `json:"date"`
//				DurationSeconds  int64    `json:"duration_seconds"`
//				Edited           string   `json:"edited"`
//				File             string   `json:"file"`
//				ForwardedFrom    string   `json:"forwarded_from"`
//				From             string   `json:"from"`
//				FromID           int64    `json:"from_id"`
//				Height           int64    `json:"height"`
//				ID               int64    `json:"id"`
//				Inviter          string   `json:"inviter"`
//				MediaType        string   `json:"media_type"`
//				Members          []string `json:"members"`
//				MimeType         string   `json:"mime_type"`
//				Photo            string   `json:"photo"`
//				ReplyToMessageID int64    `json:"reply_to_message_id"`
//				StickerEmoji     string   `json:"sticker_emoji"`
//				Text             string   `json:"text"`
//				Thumbnail        string   `json:"thumbnail"`
//				Type             string   `json:"type"`
//				Width            int64    `json:"width"`
//			} `json:"messages"`
//			Name string `json:"name"`
//			Type string `json:"type"`
//		} `json:"list"`
//	} `json:"left_chats"`
//	OtherData struct {
//		AboutMeta       string        `json:"about_meta"`
//		ChangesLog      []interface{} `json:"changes_log"`
//		CreatedStickers []struct {
//			URL string `json:"url"`
//		} `json:"created_stickers"`
//		Drafts []struct {
//			Chat     string `json:"chat"`
//			ChatName string `json:"chat_name"`
//			HTML     string `json:"html"`
//		} `json:"drafts"`
//		DraftsAbout       string `json:"drafts_about"`
//		Help              string `json:"help"`
//		InstalledStickers []struct {
//			URL string `json:"url"`
//		} `json:"installed_stickers"`
//		Ips []struct {
//			IP string `json:"ip"`
//		} `json:"ips"`
//	} `json:"other_data"`
//	PersonalInformation struct {
//		Bio         string `json:"bio"`
//		FirstName   string `json:"first_name"`
//		LastName    string `json:"last_name"`
//		PhoneNumber string `json:"phone_number"`
//		UserID      int64  `json:"user_id"`
//		Username    string `json:"username"`
//	} `json:"personal_information"`
//	ProfilePictures []struct {
//		Date  string `json:"date"`
//		Photo string `json:"photo"`
//	} `json:"profile_pictures"`
//}
