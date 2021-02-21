package instagram

import (
	"fmt"
	"gorm.io/gorm"
	"regexp"
)

var hashtagMentionsRegexp = regexp.MustCompile(fmt.Sprintf(
	"#([^{%s}]+)",
	regexp.QuoteMeta(`\"$%&'()*+,-./:;<=>?[\]^`+"`"+`{|}~#@ `)+`\n`,
))

type Hashtag struct {
	gorm.Model
	Text string `gorm:"unique"`
}

func (Hashtag) TableName() string {
	return tablePrefix + "hashtags"
}

func (h Hashtag) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"text": h.Text,
	}
}

type HashtagMention struct {
	Hashtag string
	FromIdx int
	ToIdx   int
}

func extractHashtagMentionsFromText(text string) []HashtagMention {
	var hashtagMentions []HashtagMention

	matches := hashtagMentionsRegexp.FindAllSubmatchIndex([]byte(text), -1)

	for _, bounds := range matches {
		if len(bounds) != 4 {
			continue
		}

		hashtagMentions = append(hashtagMentions, HashtagMention{
			Hashtag: text[bounds[2]:bounds[3]],
			FromIdx: bounds[0],
			ToIdx:   bounds[1],
		})
	}

	return hashtagMentions
}
