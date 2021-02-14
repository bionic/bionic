package instagram

import (
	"fmt"
	"gorm.io/gorm"
	"regexp"
)

var hashtagMentionsRegexp = regexp.MustCompile(fmt.Sprintf(
	"#([^{%s}]+)",
	regexp.QuoteMeta(`\"$%&'()*+,-./:;<=>?[\]^`+"`"+`{|}~\n#@ `),
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
	gorm.Model
	ParentID   uint   `gorm:"uniqueIndex:instagram_hashtag_mentions_key"`
	ParentType string `gorm:"uniqueIndex:instagram_hashtag_mentions_key"`
	HashtagID  uint   `gorm:"uniqueIndex:instagram_hashtag_mentions_key"`
	Hashtag    Hashtag
	FromIdx    int `gorm:"uniqueIndex:instagram_hashtag_mentions_key"`
	ToIdx      int `gorm:"uniqueIndex:instagram_hashtag_mentions_key"`
}

func (HashtagMention) TableName() string {
	return tablePrefix + "hashtag_mentions"
}

func (hm HashtagMention) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"parent_id":   hm.ParentID,
		"parent_type": hm.ParentType,
		"hashtag_id":  hm.Hashtag.ID,
		"from_idx":    hm.FromIdx,
		"to_idx":      hm.ToIdx,
	}
}

func extractHashtagMentionsFromText(text string) []HashtagMention {
	var hashtagMentions []HashtagMention

	matches := hashtagMentionsRegexp.FindAllSubmatchIndex([]byte(text), -1)

	for _, bounds := range matches {
		if len(bounds) != 4 {
			continue
		}

		hashtagMentions = append(hashtagMentions, HashtagMention{
			Hashtag: Hashtag{
				Text: text[bounds[2]:bounds[3]],
			},
			FromIdx: bounds[0],
			ToIdx:   bounds[1],
		})
	}

	return hashtagMentions
}
