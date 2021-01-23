package netflix

import (
	"github.com/BionicTeam/bionic/types"
	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type SubscriptionHistoryItem struct {
	gorm.Model
	SubscriptionOpenedTime            types.DateTime  `csv:"Subscription Opened Ts" gorm:"unique"`
	IsFreeTrial                       bool            `csv:"Is Free Trial At Signup"`
	SubscriptionClosedTime            *types.DateTime `csv:"Subscription Closed Ts"`
	IsCustomerInitiatedCancel         bool            `csv:"Is Customer Initiated Cancel"`
	SignupPlanCategory                string          `csv:"Signup Plan Category"`
	SignupMaxConcurrentStreams        int             `csv:"Signup Max Concurrent Streams"`
	SignupMaxStreamingQuality         string          `csv:"Signup Max Streaming Quality"`
	PlanChangeDate                    types.DateTime  `csv:"Plan Change Date"`
	PlanChangeOldCategory             string          `csv:"Plan Change Old Category"`
	PlanChangeOldMaxConcurrentStreams int             `csv:"Plan Change Old Max Concurrent Streams"`
	PlanChangeOldMaxStreamingQuality  string          `csv:"Plan Change Old Max Streaming Quality"`
	PlanChangeNewCategory             string          `csv:"Plan Change New Category"`
	PlanChangeNewMaxConcurrentStreams int             `csv:"Plan Change New Max Concurrent Streams"`
	PlanChangeNewMaxStreamingQuality  string          `csv:"Plan Change New Max Streaming Quality"`
	CancellationReason                string          `csv:"Cancellation Reason"`
}

func (r SubscriptionHistoryItem) TableName() string {
	return "netflix_subscription_history"
}

func (p *netflix) importSubscriptionsHistory(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var items []SubscriptionHistoryItem

	if err := gocsv.UnmarshalFile(file, &items); err != nil {
		return err
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(items, 100).
		Error
	if err != nil {
		return err
	}

	return nil
}
