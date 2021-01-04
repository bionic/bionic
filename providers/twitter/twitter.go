package twitter

import (
	"github.com/shekhirin/bionic-cli/providers/provider"
	"gorm.io/gorm"
)

type twitter struct {
	db *gorm.DB
}

func New(db *gorm.DB) provider.Provider {
	return &twitter{
		db: db,
	}
}

func (p *twitter) Models() []interface{} {
	return []interface{}{
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

	if err := p.processLikes(inputPath); err != nil {
		return err
	}

	if err := p.processDirectMessages(inputPath); err != nil {
		return err
	}

	if err := p.processTweets(inputPath); err != nil {
		return err
	}

	return nil
}
