package ofx

import (
	"github.com/bionic-dev/bionic/internal/provider/database"
	_ "github.com/bionic-dev/bionic/testinit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestOFX_importStatement(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:?_loc=UTC"), &gorm.Config{})
	require.NoError(t, err)

	p := OFX{Database: database.New(db)}
	require.NoError(t, p.Migrate())

	require.NoError(t, p.importStatement("testdata/ofx/statement.ofx"))

	var accounts []Account
	require.NoError(t, p.DB().Preload("Transactions").Find(&accounts).Error)
	require.Len(t, accounts, 1)

	expectedAccount := Account{
		Currency:  "RUB",
		BankID:    "Tinkoff",
		AccountID: "123",
		Type:      "CHECKING",
		Transactions: []Transaction{
			{
				Type:         "CREDIT",
				Time:         DateTime{Time: time.Date(2021, 3, 18, 0, 0, 0, 0, time.UTC)},
				Amount:       -380.25,
				FITID:        "eb93cd****************2c07ae65",
				Name:         "",
				Memo:         "Bools: Gumroad.Co* Daniel Vas, 5192, *****",
				Currency:     "",
				CurrencyRate: 0,
			},
		},
	}
	assertAccount(t, expectedAccount, accounts[0])

	require.NoError(t, p.importStatement("testdata/ofx/statement.ofx"))

	var newAccounts []Account
	require.NoError(t, p.DB().Preload("Transactions").Find(&newAccounts).Error)
	require.Len(t, newAccounts, 1)
	assertAccount(t, accounts[0], newAccounts[0])
}

func assertAccount(t *testing.T, expected, actual Account) {
	expected = convertAccount(expected)
	actual = convertAccount(actual)
	assert.EqualValues(t, expected, actual)
}

func convertAccount(account Account) Account {
	account.Model = gorm.Model{}
	for i := range account.Transactions {
		account.Transactions[i].Model = gorm.Model{}
		account.Transactions[i].AccountID = 0
		account.Transactions[i].Account = Account{}
	}
	return account
}
