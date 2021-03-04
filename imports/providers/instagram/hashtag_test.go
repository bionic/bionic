package instagram

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractHashtagMentionsFromText(t *testing.T) {
	text := "#hashtag #hash#tag # hashtag #hash$tag #1tag"

	hashtagMentions := extractHashtagMentionsFromText(text)

	assert.Equal(t, []HashtagMention{
		{
			Hashtag: "hashtag",
			FromIdx: 0,
			ToIdx:   8,
		},
		{
			Hashtag: "hash",
			FromIdx: 9,
			ToIdx:   14,
		},
		{
			Hashtag: "tag",
			FromIdx: 14,
			ToIdx:   18,
		},
		{
			Hashtag: "hash",
			FromIdx: 29,
			ToIdx:   34,
		},
		{
			Hashtag: "1tag",
			FromIdx: 39,
			ToIdx:   44,
		},
	}, hashtagMentions)
}
