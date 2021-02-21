package instagram

import (
	"gorm.io/gorm"
	"regexp"
)

var userMentionsRegexp = regexp.MustCompile(`(?:^|[^\w])(@([\w_.]+))`)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
}

func (User) TableName() string {
	return tablePrefix + "users"
}

func (u User) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"username": u.Username,
	}
}

type UserMention struct {
	Username string
	FromIdx  int
	ToIdx    int
}

func extractUserMentionsFromText(text string) []UserMention {
	var userMentions []UserMention

	matches := userMentionsRegexp.FindAllSubmatchIndex([]byte(text), -1)

	for _, bounds := range matches {
		if len(bounds) != 6 {
			continue
		}

		userMentions = append(userMentions, UserMention{
			Username: text[bounds[4]:bounds[5]],
			FromIdx:  bounds[2],
			ToIdx:    bounds[3],
		})
	}

	return userMentions
}
