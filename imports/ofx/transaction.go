package ofx

import (
	"github.com/aclindsa/xml"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Type         string   `xml:"TRNTYPE"`
	Time         DateTime `xml:"DTPOSTED"`
	Amount       float32  `xml:"TRNAMT"`
	FITID        string   `xml:"FITID"`
	Name         string   `xml:"NAME"`
	Memo         string   `xml:"MEMO"`
	Currency     string
	CurrencyRate float32
	AccountID    uint
	Account      Account
}

func (Transaction) TableName() string {
	return tablePrefix + "transactions"
}

func (t *Transaction) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	type Alias Transaction

	var data struct {
		Alias
		Currency struct {
			Symbol string  `xml:"CURSYM"`
			Rate   float32 `xml:"CURRATE"`
		} `xml:"CURRENCY"`
	}

	if err := decoder.DecodeElement(&data, &start); err != nil {
		return err
	}

	*t = Transaction(data.Alias)
	t.Currency = data.Currency.Symbol
	t.CurrencyRate = data.Currency.Rate

	return nil
}

func (t Transaction) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"time":   t.Time,
		"amount": t.Amount,
		"fit_id": t.FITID,
	}
}
