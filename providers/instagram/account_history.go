package instagram

import (
	"encoding/json"
	"github.com/bionic-dev/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
)

type AccountHistoryAction string

const (
	AccountHistoryLogin  AccountHistoryAction = "login"
	AccountHistoryLogout AccountHistoryAction = "logout"
)

type AccountHistoryItem struct {
	gorm.Model
	Action       AccountHistoryAction `gorm:"uniqueIndex:instagram_account_history_key"`
	CookieName   string               `json:"cookie_name"`
	IPAddress    string               `json:"ip_address" gorm:"uniqueIndex:instagram_account_history_key"`
	LanguageCode string               `json:"language_code"`
	Timestamp    types.DateTime       `json:"timestamp" gorm:"uniqueIndex:instagram_account_history_key"`
	UserAgent    string               `json:"user_agent"`
	DeviceID     *string              `json:"device_id"`
}

func (AccountHistoryItem) TableName() string {
	return tablePrefix + "account_history"
}

func (ah *AccountHistoryItem) UnmarshalJSON(b []byte) error {
	type Alias AccountHistoryItem

	var data struct {
		Alias
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*ah = AccountHistoryItem(data.Alias)

	if ah.DeviceID != nil && *ah.DeviceID == "" {
		ah.DeviceID = nil
	}

	return nil
}

type RegistrationInfo struct {
	gorm.Model
	RegistrationUsername    string         `json:"registration_username" gorm:"unique"`
	IPAddress               string         `json:"ip_address"`
	RegistrationTime        types.DateTime `json:"registration_time"`
	RegistrationEmail       string         `json:"registration_email"`
	RegistrationPhoneNumber *string        `json:"registration_phone_number"`
	DeviceName              *string        `json:"device_name"`
}

func (RegistrationInfo) TableName() string {
	return tablePrefix + "registration_info"
}

func (ri *RegistrationInfo) UnmarshalJSON(b []byte) error {
	type Alias RegistrationInfo

	var data struct {
		Alias
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*ri = RegistrationInfo(data.Alias)

	if ri.RegistrationPhoneNumber != nil && *ri.RegistrationPhoneNumber == "" {
		ri.RegistrationPhoneNumber = nil
	}
	if ri.DeviceName != nil && *ri.DeviceName == "" {
		ri.DeviceName = nil
	}

	return nil
}

func (p *instagram) importAccountHistory(inputPath string) error {
	var data struct {
		LoginHistory     []AccountHistoryItem `json:"login_history"`
		RegistrationInfo RegistrationInfo     `json:"registration_info"`
		LogoutHistory    []AccountHistoryItem `json:"logout_history"`
	}

	bytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	for i := range data.LoginHistory {
		data.LoginHistory[i].Action = AccountHistoryLogin
	}
	for i := range data.LogoutHistory {
		data.LogoutHistory[i].Action = AccountHistoryLogout
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(&data.RegistrationInfo).
		Error
	if err != nil {
		return err
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(&data.LoginHistory, 100).
		Error
	if err != nil {
		return err
	}

	err = p.DB().
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		CreateInBatches(&data.LogoutHistory, 100).
		Error
	if err != nil {
		return err
	}

	return nil
}
