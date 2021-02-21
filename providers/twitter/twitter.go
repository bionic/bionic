package twitter

import (
	"encoding/json"
	"github.com/bionic-dev/bionic/providers/provider"
	"gorm.io/gorm"
	"io/ioutil"
	"path"
	"strings"
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

func (p *twitter) Migrate() error {
	return p.DB().AutoMigrate(
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
		&AgeInfoRecord{},
		&AdImpression{},
		&DeviceInfo{},
		&TargetingCriterion{},
		&Account{},
		&ScreenNameChange{},
		&EmailAddressChange{},
		&LoginIP{},
	)
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
		provider.NewImportFn(
			"Age Info",
			p.importAgeInfo,
			path.Join(inputPath, "data", "ageinfo.js"),
		),
		provider.NewImportFn(
			"Ad Impressions",
			p.importAdImpressions,
			path.Join(inputPath, "data", "ad-impressions.js"),
		),
		provider.NewImportFn(
			"Account",
			p.importAccount,
			path.Join(inputPath, "data"),
		),
	}, nil
}

func readJSON(path string, prefix string, dst interface{}) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	data := string(bytes)
	data = strings.TrimPrefix(data, prefix)

	if err := json.Unmarshal([]byte(data), &dst); err != nil {
		return err
	}

	return nil
}
