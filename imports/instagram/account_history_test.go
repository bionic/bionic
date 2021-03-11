package instagram

import (
	"github.com/bionic-dev/bionic/internal/provider/database"
	"github.com/bionic-dev/bionic/pkg/ptr"
	_ "github.com/bionic-dev/bionic/testinit"
	"github.com/bionic-dev/bionic/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestInstagram_importAccountHistory(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	p := instagram{Database: database.New(db)}
	require.NoError(t, p.Migrate())

	require.NoError(t, p.importAccountHistory("testdata/instagram/account_history.json"))

	var accountHistory []AccountHistoryItem
	require.NoError(t, db.Find(&accountHistory).Error)
	require.Len(t, accountHistory, 4)

	assertAccountHistoryItem(t, AccountHistoryItem{
		Action:       AccountHistoryLogin,
		CookieName:   "*************************Oxb",
		IPAddress:    "34.207.117.203",
		LanguageCode: "en",
		Timestamp:    types.DateTime(time.Date(2020, 10, 31, 20, 37, 21, 0, time.UTC)),
		UserAgent:    "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36",
		DeviceID:     ptr.String("SOME-DEVICE-ID"),
	}, accountHistory[0])
	assertAccountHistoryItem(t, AccountHistoryItem{
		Action:       AccountHistoryLogin,
		CookieName:   "*************************VDE",
		IPAddress:    "95.161.223.248",
		LanguageCode: "en",
		Timestamp:    types.DateTime(time.Date(2020, 10, 31, 20, 1, 39, 0, time.UTC)),
		UserAgent:    "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36",
		DeviceID:     nil,
	}, accountHistory[1])
	assertAccountHistoryItem(t, AccountHistoryItem{
		Action:       AccountHistoryLogout,
		CookieName:   "*************************FJi",
		IPAddress:    "145.255.233.150",
		LanguageCode: "en",
		Timestamp:    types.DateTime(time.Date(2019, 4, 3, 20, 55, 30, 0, time.UTC)),
		UserAgent:    "Instagram 87.0.0.16.99 (iPhone9,3; iOS 12_1_4; en_RU; en-RU; scale=2.00; gamut=wide; 750x1334; 147928430) AppleWebKit/420+",
		DeviceID:     ptr.String("SOME-DEVICE-ID"),
	}, accountHistory[2])
	assertAccountHistoryItem(t, AccountHistoryItem{
		Action:       AccountHistoryLogout,
		CookieName:   "*************************14m",
		IPAddress:    "145.255.233.150",
		LanguageCode: "en",
		Timestamp:    types.DateTime(time.Date(2019, 03, 28, 13, 51, 0, 0, time.UTC)),
		UserAgent:    "Instagram 52.0.0.8.83 (iPhone; CPU iPhone OS 11_4 like Mac OS X; en_US; en-US; scale=2.00; 750x1334) AppleWebKit/605.1.15",
		DeviceID:     nil,
	}, accountHistory[3])

	var registrationInfo []RegistrationInfo
	require.NoError(t, db.Find(&registrationInfo).Error)
	require.Len(t, registrationInfo, 1)

	assert.Equal(t, "cool_nickname", registrationInfo[0].RegistrationUsername)
	assert.Equal(t, "91.122.198.16", registrationInfo[0].IPAddress)
	assert.EqualValues(t, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), registrationInfo[0].RegistrationTime)
	assert.Equal(t, "admin@gmail.com", registrationInfo[0].RegistrationEmail)
	assert.Nil(t, registrationInfo[0].RegistrationPhoneNumber)
	assert.Nil(t, registrationInfo[0].DeviceName)
}

func assertAccountHistoryItem(t *testing.T, expected, actual AccountHistoryItem) {
	assert.Equal(t, expected.Action, actual.Action)
	assert.Equal(t, expected.CookieName, actual.CookieName)
	assert.Equal(t, expected.IPAddress, actual.IPAddress)
	assert.Equal(t, expected.LanguageCode, actual.LanguageCode)
	assert.Equal(t, expected.Timestamp, actual.Timestamp)
	assert.Equal(t, expected.UserAgent, actual.UserAgent)
	assert.Equal(t, expected.DeviceID, actual.DeviceID)
}
