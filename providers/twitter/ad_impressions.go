package twitter

import (
	"encoding/json"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdImpression struct {
	gorm.Model
	DeviceInfoID             int
	DeviceInfo               DeviceInfo `json:"deviceInfo"`
	DisplayLocation          string     `json:"displayLocation"`
	PromotedTweetID          int
	PromotedTweet            Tweet
	AdvertiserID             int
	Advertiser               Advertiser
	MatchedTargetingCriteria []TargetingCriterion `json:"matchedTargetingCriteria" gorm:"many2many:twitter_ad_impressions_matched_targeting_criteria"`
	ImpressionTime           types.DateTime       `json:"impressionTime" gorm:"unique"`
}

func (AdImpression) TableName() string {
	return tablePrefix + "ad_impressions"
}

func (ai *AdImpression) UnmarshalJSON(b []byte) error {
	type alias AdImpression

	var data struct {
		alias
		PromotedTweetInfo struct {
			TweetID   int      `json:"tweetId,string"`
			TweetText string   `json:"tweetText"`
			URLs      []string `json:"urls"`
			MediaURLs []string `json:"mediaUrls"`
		} `json:"promotedTweetInfo"`
		AdvertiserInfo struct {
			AdvertiserName string `json:"advertiserName"`
			ScreenName     string `json:"screenName"`
		} `json:"advertiserInfo"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*ai = AdImpression(data.alias)

	ai.PromotedTweet.ID = data.PromotedTweetInfo.TweetID
	ai.PromotedTweet.FullText = data.PromotedTweetInfo.TweetText

	for _, url := range data.PromotedTweetInfo.URLs {
		ai.PromotedTweet.Entities.URLs = append(ai.PromotedTweet.Entities.URLs, TweetURL{
			URL: URL{
				URL: url,
			},
		})
	}

	for _, url := range data.PromotedTweetInfo.MediaURLs {
		ai.PromotedTweet.Entities.Media = append(ai.PromotedTweet.Entities.Media, TweetMedia{
			URL: url,
		})
	}

	ai.Advertiser.Name = data.AdvertiserInfo.ScreenName

	return nil
}

type DeviceInfo struct {
	gorm.Model
	Identifier string `json:"deviceId" gorm:"uniqueIndex:twitter_device_infos_key"`
	Type       string `json:"deviceType" gorm:"uniqueIndex:twitter_device_infos_key"`
	OsType     string `json:"osType" gorm:"uniqueIndex:twitter_device_infos_key"`
}

func (DeviceInfo) TableName() string {
	return tablePrefix + "device_infos"
}

type TargetingCriterion struct {
	gorm.Model
	TargetingType  string `json:"targetingType" gorm:"uniqueIndex:twitter_targeting_criterion_key"`
	TargetingValue string `json:"targetingValue" gorm:"uniqueIndex:twitter_targeting_criterion_key"`
}

func (TargetingCriterion) TableName() string {
	return tablePrefix + "targeting_criteria"
}

func (p *twitter) importAdImpressions(inputPath string) error {
	var ads []struct {
		Ad struct {
			AdsUserData struct {
				AdImpressions struct {
					Impressions []AdImpression `json:"impressions"`
				} `json:"adImpressions"`
			} `json:"adsUserData"`
		} `json:"ad"`
	}

	if err := readJSON(
		inputPath,
		"window.YTD.ad_impressions.part0 = ",
		&ads,
	); err != nil {
		return err
	}

	var adImpressions []AdImpression

	for i, ad := range ads {
		for j := range ad.Ad.AdsUserData.AdImpressions.Impressions {
			adImpression := &ads[i].Ad.AdsUserData.AdImpressions.Impressions[j]

			err := p.DB().
				FirstOrCreate(&adImpression.DeviceInfo, map[string]interface{}{
					"identifier": adImpression.DeviceInfo.Identifier,
					"type":       adImpression.DeviceInfo.Type,
					"os_type":    adImpression.DeviceInfo.OsType,
				}).
				Error
			if err != nil {
				return err
			}

			err = p.DB().
				Find(&adImpression.PromotedTweet.Entities, map[string]interface{}{
					"tweet_id": adImpression.PromotedTweet.ID,
				}).
				Error
			if err != nil {
				return err
			}

			for k := range adImpression.PromotedTweet.Entities.URLs {
				url := &adImpression.PromotedTweet.Entities.URLs[k]

				err = p.DB().
					FirstOrCreate(&url.URL, map[string]interface{}{
						"url": url.URL.URL,
					}).
					Error
				if err != nil {
					return err
				}

				err = p.DB().
					FirstOrCreate(url, map[string]interface{}{
						"tweet_entities_id": adImpression.PromotedTweet.Entities.ID,
						"url_id":            url.URL.ID,
					}).
					Error
				if err != nil {
					return err
				}
			}

			for k := range adImpression.PromotedTweet.Entities.Media {
				media := &adImpression.PromotedTweet.Entities.Media[k]

				err = p.DB().
					FirstOrCreate(media, map[string]interface{}{
						"url":               media.URL,
						"tweet_entities_id": adImpression.PromotedTweet.Entities.ID,
					}).
					Error
				if err != nil {
					return err
				}
			}

			err = p.DB().
				FirstOrCreate(&adImpression.PromotedTweet.Entities, map[string]interface{}{
					"tweet_id": adImpression.PromotedTweet.ID,
				}).
				Error
			if err != nil {
				return err
			}

			err = p.DB().
				FirstOrCreate(&adImpression.PromotedTweet, map[string]interface{}{
					"id": adImpression.PromotedTweet.ID,
				}).
				Error
			if err != nil {
				return err
			}

			err = p.DB().
				FirstOrCreate(&adImpression.Advertiser, map[string]interface{}{
					"name": adImpression.Advertiser.Name,
				}).
				Error
			if err != nil {
				return err
			}

			for k := range adImpression.MatchedTargetingCriteria {
				criterion := &adImpression.MatchedTargetingCriteria[k]

				err = p.DB().
					FirstOrCreate(criterion, map[string]interface{}{
						"targeting_type":  criterion.TargetingType,
						"targeting_value": criterion.TargetingValue,
					}).
					Error
				if err != nil {
					return err
				}
			}

			adImpressions = append(adImpressions, *adImpression)
		}
	}

	err := p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(adImpressions, 100).
		Error
	if err != nil {
		return err
	}

	return nil
}
