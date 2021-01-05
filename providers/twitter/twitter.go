package twitter

import (
	"github.com/shekhirin/bionic-cli/providers/provider"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"path"
)

type twitter struct {
	provider.Database
}

func New(db *gorm.DB) provider.Provider {
	return &twitter{
		Database: provider.NewDatabase(db),
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

func (p *twitter) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	if !provider.IsPathDir(inputPath) {
		return nil, provider.ErrInputPathShouldBeDirectory
	}

	return []provider.ImportFn{
		{
			Fn:        p.importLikes,
			InputPath: path.Join(inputPath, "data", "like.js"),
		},
		{
			Fn:        p.importDirectMessages,
			InputPath: path.Join(inputPath, "data", "direct-messages.js"),
		},
		{
			Fn:        p.importTweets,
			InputPath: path.Join(inputPath, "data", "tweet.js"),
		},
	}, nil
}
