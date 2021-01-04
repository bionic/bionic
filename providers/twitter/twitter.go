package twitter

import (
	"github.com/shekhirin/bionic-cli/providers/provider"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"path"
)

type twitter struct {
	db *gorm.DB
}

func New(db *gorm.DB) provider.Provider {
	return &twitter{
		db: db,
	}
}

func (p *twitter) Models() []schema.Tabler {
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
	}
}

func (p *twitter) Process(inputPath string) error {
	if !provider.IsPathDir(inputPath) {
		return provider.ErrInputPathShouldBeDirectory
	}

	if err := p.processLikes(path.Join(inputPath, "data", "like.js")); err != nil {
		return err
	}

	if err := p.processDirectMessages(path.Join(inputPath, "data", "direct-messages.js")); err != nil {
		return err
	}

	if err := p.processTweets(path.Join(inputPath, "data", "tweet.js")); err != nil {
		return err
	}

	return nil
}
