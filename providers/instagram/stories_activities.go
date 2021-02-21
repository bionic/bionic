package instagram

import (
	"encoding/json"
	"errors"
	"github.com/bionic-dev/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
)

type StoryActivityType string

const (
	StoryActivityPoll        StoryActivityType = "poll"
	StoryActivityEmojiSlider StoryActivityType = "emoji_slider"
	StoryActivityQuestion    StoryActivityType = "question"
	StoryActivityCountdown   StoryActivityType = "countdown"
	StoryActivityQuiz        StoryActivityType = "quiz"
)

type StoryActivityItem struct {
	gorm.Model
	Type      StoryActivityType `gorm:"uniqueIndex:instagram_story_activities_key"`
	UserID    int               `gorm:"uniqueIndex:instagram_story_activities_key"`
	User      User
	Timestamp types.DateTime `gorm:"uniqueIndex:instagram_story_activities_key"`
}

func (StoryActivityItem) TableName() string {
	return tablePrefix + "story_activities"
}

func (sai *StoryActivityItem) UnmarshalJSON(b []byte) error {
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
		Polls        []StoryActivityItem `json:"polls"`
		EmojiSliders []StoryActivityItem `json:"emoji_sliders"`
		Questions    []StoryActivityItem `json:"questions"`
		Countdowns   []StoryActivityItem `json:"countdowns"`
		Quizzes      []StoryActivityItem `json:"quizzes"`
	}

	bytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	if err := p.insertStoryActivities(data.Polls, StoryActivityPoll); err != nil {
		return err
	}

	if err := p.insertStoryActivities(data.EmojiSliders, StoryActivityEmojiSlider); err != nil {
		return err
	}

	if err := p.insertStoryActivities(data.Questions, StoryActivityQuestion); err != nil {
		return err
	}

	if err := p.insertStoryActivities(data.Countdowns, StoryActivityCountdown); err != nil {
		return err
	}

	if err := p.insertStoryActivities(data.Quizzes, StoryActivityQuiz); err != nil {
		return err
	}

	return nil
}

func (p *instagram) insertStoryActivities(activities []StoryActivityItem, activityType StoryActivityType) error {
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
