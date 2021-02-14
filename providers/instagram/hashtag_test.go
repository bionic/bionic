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
			Hashtag: Hashtag{Text: "hashtag"},
			FromIdx: 0,
			ToIdx: 8,
		},
		{
			Hashtag: Hashtag{Text: "hash"},
			FromIdx: 9,
			ToIdx: 14,
		},
		{
			Hashtag: Hashtag{Text: "tag"},
			FromIdx: 14,
			ToIdx: 18,
		},
		{
			Hashtag: Hashtag{Text: "hash"},
			FromIdx: 29,
			ToIdx: 34,
		},
		{
			Hashtag: Hashtag{Text: "1tag"},
			FromIdx: 39,
			ToIdx: 44,
		},
	}, hashtagMentions)
}
