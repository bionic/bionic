package instagram

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractUserMentionsFromText(t *testing.T) {
	text := "@beginning This is a @with.dot of some cool features that @under_score be useful but don't. " +
		"look at this email@address.ignored @mention! I like #nylas but I don't like to go to this apple.com?a#url. " +
		"I also don't like the ### comment blocks. But #msft is cool."

	userMentions := extractUserMentionsFromText(text)

	assert.Equal(t, []UserMention{
		{
			Username: "beginning",
			FromIdx:  0,
			ToIdx:    10,
		},
		{
			Username: "with.dot",
			FromIdx:  21,
			ToIdx:    30,
		},
		{
			Username: "under_score",
			FromIdx:  58,
			ToIdx:    70,
		},
		{
			Username: "mention",
			FromIdx:  127,
			ToIdx:    135,
		},
	}, userMentions)
}
