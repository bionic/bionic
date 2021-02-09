package netflix

import (
	"github.com/BionicTeam/bionic/types"
	"github.com/gocarina/gocsv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type BillingHistoryItem struct {
	gorm.Model
	TransactionDate        types.DateTime `csv:"Transaction Date" gorm:"uniqueIndex:netflix_billing_history_key"`
	ServicePeriodStartDate types.DateTime `csv:"Service Period Start Date"`
	ServicePeriodEndDate   types.DateTime `csv:"Service Period End Date"`
	Description            string         `csv:"Description"`
	PaymentType            string         `csv:"Payment Type" gorm:"uniqueIndex:netflix_billing_history_key"`
	MopLast4               string         `csv:"Mop Last 4" gorm:"uniqueIndex:netflix_billing_history_key"`
	MopCreationDate        types.DateTime `csv:"Mop Creation Date"`
	MopPmtProcessorDesc    string         `csv:"Mop Pmt Processor Desc"`
	ItemPriceAmt           float32        `csv:"Item Price Amt"`
	Currency               string         `csv:"Currency"`
	TaxAmt                 float32        `csv:"Tax Amt"`
	GrossSaleAmt           float32        `csv:"Gross Sale Amt"`
	PmtTxnType             string         `csv:"Pmt Txn Type" gorm:"uniqueIndex:netflix_billing_history_key"`
	PmtStatus              string         `csv:"Pmt Status" gorm:"uniqueIndex:netflix_billing_history_key"`
	FinalInvoiceResult     string         `csv:"Final Invoice Result" gorm:"uniqueIndex:netflix_billing_history_key"`
	Country                string         `csv:"Country"`
	NextBillingDate        types.DateTime `csv:"Next Billing Date"`
}

func (r BillingHistoryItem) TableName() string {
	return tablePrefix + "billing_history"
}

func (p *netflix) importBillingHistory(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	var items []BillingHistoryItem

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
