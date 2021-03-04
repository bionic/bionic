package instagram

import (
	"encoding/json"
	"errors"
	"github.com/bionic-dev/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
)

type StoriesActivityType string

const (
	StoriesActivityPoll        StoriesActivityType = "poll"
	StoriesActivityEmojiSlider StoriesActivityType = "emoji_slider"
	StoriesActivityQuestion    StoriesActivityType = "question"
	StoriesActivityCountdown   StoriesActivityType = "countdown"
	StoriesActivityQuiz        StoriesActivityType = "quiz"
)

type StoriesActivityItem struct {
	gorm.Model
	Type      StoriesActivityType `gorm:"uniqueIndex:instagram_stories_activities_key"`
	UserID    int                 `gorm:"uniqueIndex:instagram_stories_activities_key"`
	User      User
	Timestamp types.DateTime `gorm:"uniqueIndex:instagram_stories_activities_key"`
}

func (StoriesActivityItem) TableName() string {
	return tablePrefix + "stories_activities"
}

func (sai *StoriesActivityItem) UnmarshalJSON(b []byte) error {
	var data []string

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if len(data) != 2 {
		return errors.New("incorrect stories activities format")
	}

	if err := sai.Timestamp.UnmarshalText([]byte(data[0])); err != nil {
		return err
	}

	sai.User.Username = data[1]

	return nil
}

func (p *instagram) importStoriesActivities(inputPath string) error {
	var data struct {
		Polls        []StoriesActivityItem `json:"polls"`
		EmojiSliders []StoriesActivityItem `json:"emoji_sliders"`
		Questions    []StoriesActivityItem `json:"questions"`
		Countdowns   []StoriesActivityItem `json:"countdowns"`
		Quizzes      []StoriesActivityItem `json:"quizzes"`
	}

	bytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	if err := p.insertStoriesActivities(data.Polls, StoriesActivityPoll); err != nil {
		return err
	}

	if err := p.insertStoriesActivities(data.EmojiSliders, StoriesActivityEmojiSlider); err != nil {
		return err
	}

	if err := p.insertStoriesActivities(data.Questions, StoriesActivityQuestion); err != nil {
		return err
	}

	if err := p.insertStoriesActivities(data.Countdowns, StoriesActivityCountdown); err != nil {
		return err
	}

	if err := p.insertStoriesActivities(data.Quizzes, StoriesActivityQuiz); err != nil {
		return err
	}

	return nil
}

func (p *instagram) insertStoriesActivities(activities []StoriesActivityItem, activityType StoriesActivityType) error {
	for i := range activities {
		activities[i].Type = activityType

		err := p.DB().
			FirstOrCreate(&activities[i].User, activities[i].User.Conditions()).
			Error
		if err != nil {
			return err
		}
	}

	return p.DB().
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(activities, 1000).
		Error
}
