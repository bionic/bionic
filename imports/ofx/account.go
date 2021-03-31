package ofx

import (
	"github.com/aclindsa/xml"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Currency     string
	BankID       string
	AccountID    string
	Type         string
	Transactions []Transaction
}

func (Account) TableName() string {
	return tablePrefix + "accounts"
}

func (a *Account) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var data struct {
		AccountInfo struct {
			Currency string `xml:"CURDEF"`
			Meta     struct {
				BankID    string `xml:"BANKID"`
				AccountID string `xml:"ACCTID"`
				Type      string `xml:"ACCTTYPE"`
			} `xml:"BANKACCTFROM"`
			TransactionsInfo struct {
				Transactions []Transaction `xml:"STMTTRN"`
			} `xml:"BANKTRANLIST"`
		} `xml:"STMTRS"`
	}

	if err := decoder.DecodeElement(&data, &start); err != nil {
		return err
	}

	a.Currency = data.AccountInfo.Currency
	a.BankID = data.AccountInfo.Meta.BankID
	a.AccountID = data.AccountInfo.Meta.AccountID
	a.Type = data.AccountInfo.Meta.Type
	a.Transactions = data.AccountInfo.TransactionsInfo.Transactions

	return nil
}

func (a Account) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"bank_id":    a.BankID,
		"account_id": a.AccountID,
	}
}
