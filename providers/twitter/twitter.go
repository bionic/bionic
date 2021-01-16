package twitter

import (
	"github.com/shekhirin/bionic-cli/providers/provider"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"path"
)

const name = "twitter"
const tablePrefix = "twitter_"

type twitter struct {
	provider.Database
}

func New(db *gorm.DB) provider.Provider {
	return &twitter{
		Database: provider.NewDatabase(db),
	}
}

func (twitter) Name() string {
	return name
}

func (twitter) TablePrefix() string {
	return tablePrefix
}

func (twitter) Models() []schema.Tabler {
	return []schema.Tabler{
		&Like{},
		&URL{},
		&Conversation{},
		&DirectMessage{},
		&DirectMessageReaction{},
		&Tweet{},
		&TweetEntities{},
		&TweetHashtag{},
		&TweetMedia{},
		&TweetUserMention{},
		&TweetURL{},
		&PersonalizationRecord{},
		&LanguageRecord{},
		&GenderInfo{},
		&InterestRecord{},
		&AudienceAndAdvertiserRecord{},
		&Advertiser{},
		&Show{},
		&Location{},
		&InferredAgeInfoRecord{},
	}
}

func (p *twitter) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	if !provider.IsPathDir(inputPath) {
		return nil, provider.ErrInputPathShouldBeDirectory
	}

	return []provider.ImportFn{
		provider.NewImportFn(
			"Likes",
			p.importLikes,
			path.Join(inputPath, "data", "like.js"),
		),
		provider.NewImportFn(
			"Direct Messages",
			p.importDirectMessages,
			path.Join(inputPath, "data", "direct-messages.js"),
		),
		provider.NewImportFn(
			"Tweets",
			p.importTweets,
			path.Join(inputPath, "data", "tweet.js"),
		),
		provider.NewImportFn(
			"Personalization",
			p.importPersonalization,
			path.Join(inputPath, "data", "personalization.js"),
		),
	}, nil
}
