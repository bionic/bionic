package ofx

import (
	"github.com/aclindsa/xml"
	"io/ioutil"
)

func (p *OFX) importStatement(filePath string) error {
	var data struct {
		Info struct {
			Accounts []Account `xml:"STMTTRNRS"`
		} `xml:"BANKMSGSRSV1"`
	}

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := xml.Unmarshal(bytes, &data); err != nil {
		return err
	}

	for _, account := range data.Info.Accounts {
		transactions := account.Transactions
		account.Transactions = nil

		err = p.DB().FirstOrCreate(&account, account.Conditions()).Error
		if err != nil {
			return err
		}

		for _, transaction := range transactions {
			transaction.AccountID = account.ID
			err = p.DB().FirstOrCreate(&transaction, transaction.Conditions()).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
