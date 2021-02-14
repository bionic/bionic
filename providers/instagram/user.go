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
	gorm.Model
	ParentID   uint   `gorm:"uniqueIndex:instagram_user_mentions_key"`
	ParentType string `gorm:"uniqueIndex:instagram_user_mentions_key"`
	UserID     uint   `gorm:"uniqueIndex:instagram_user_mentions_key"`
	User       User
	FromIdx    int `gorm:"uniqueIndex:instagram_user_mentions_key"`
	ToIdx      int `gorm:"uniqueIndex:instagram_user_mentions_key"`
}

func (UserMention) TableName() string {
	return tablePrefix + "user_mentions"
}

func (um UserMention) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"parent_id":   um.ParentID,
		"parent_type": um.ParentType,
		"user_id":     um.User.ID,
		"from_idx":    um.FromIdx,
		"to_idx":      um.ToIdx,
	}
}

func extractUserMentionsFromText(text string) []UserMention {
	var userMentions []UserMention

	matches := userMentionsRegexp.FindAllSubmatchIndex([]byte(text), -1)

	for _, bounds := range matches {
		if len(bounds) != 6 {
			continue
		}

		userMentions = append(userMentions, UserMention{
			User: User{
				Username: text[bounds[4]:bounds[5]],
			},
			FromIdx: bounds[2],
			ToIdx:   bounds[3],
		})
	}

	return userMentions
}
