package twitter

import (
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"path"
)

type accountCreationIP struct {
	CreationIP string `json:"userCreationIp"`
}

type accountTimezone struct {
	Timezone string `json:"timeZone"`
}

type accountVerified struct {
	Verified bool `json:"verified"`
}

type Account struct {
	gorm.Model
	accountCreationIP
	accountTimezone
	accountVerified
	ID                  int            `json:"accountId,string"`
	Email               string         `json:"email"`
	CreatedVia          string         `json:"createdVia"`
	Username            string         `json:"username"`
	Created             types.DateTime `json:"createdAt"`
	DisplayName         string         `json:"accountDisplayName"`
	ScreenNameChanges   []ScreenNameChange
	EmailAddressChanges []EmailAddressChange
	LoginIPs            []LoginIP
}

func (Account) TableName() string {
	return tablePrefix + "accounts"
}

type ScreenNameChange struct {
	gorm.Model
	AccountID   int            `gorm:"uniqueIndex:twitter_screen_name_changes_key"`
	ChangedAt   types.DateTime `json:"changedAt" gorm:"uniqueIndex:twitter_screen_name_changes_key"`
	ChangedFrom string         `json:"changedFrom" gorm:"uniqueIndex:twitter_screen_name_changes_key"`
	ChangedTo   *string        `json:"changedTo" gorm:"uniqueIndex:twitter_screen_name_changes_key"`
}

func (ScreenNameChange) TableName() string {
	return tablePrefix + "screen_name_changes"
}

func (snc ScreenNameChange) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"account_id":   snc.AccountID,
		"changed_at":   snc.ChangedAt,
		"changed_from": snc.ChangedFrom,
		"changed_to":   snc.ChangedTo,
	}
}

type EmailAddressChange struct {
	gorm.Model
	AccountID   int            `gorm:"uniqueIndex:twitter_email_address_changes_key"`
	ChangedAt   types.DateTime `json:"changedAt" gorm:"uniqueIndex:twitter_email_address_changes_key"`
	ChangedFrom string         `json:"changedFrom" gorm:"uniqueIndex:twitter_email_address_changes_key"`
	ChangedTo   *string        `json:"changedTo" gorm:"uniqueIndex:twitter_email_address_changes_key"`
}

func (EmailAddressChange) TableName() string {
	return tablePrefix + "email_address_changes"
}

func (eac EmailAddressChange) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"account_id":   eac.AccountID,
		"changed_at":   eac.ChangedAt,
		"changed_from": eac.ChangedFrom,
		"changed_to":   eac.ChangedTo,
	}
}

type LoginIP struct {
	gorm.Model
	AccountID int            `json:"accountId,string" gorm:"uniqueIndex:twitter_login_ips_key"`
	IP        string         `json:"loginIp" gorm:"uniqueIndex:twitter_login_ips_key"`
	Created   types.DateTime `json:"createdAt" gorm:"uniqueIndex:twitter_login_ips_key"`
}

func (LoginIP) TableName() string {
	return tablePrefix + "login_ips"
}

func (li LoginIP) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"account_id": li.AccountID,
		"ip":         li.IP,
		"created":    li.Created,
	}
}

func (p *twitter) importAccount(inputPath string) error {
	var data = make([]struct {
		Account           Account `json:"account"`
		AccountCreationIP struct {
			accountCreationIP
			AccountId int `json:"accountId,string" gorm:"-"`
		} `json:"accountCreationIp"`
		AccountTimezone struct {
			accountTimezone
			AccountId int `json:"accountId,string" gorm:"-"`
		} `json:"accountTimezone"`
		AccountVerified struct {
			accountVerified
			AccountId int `json:"accountId,string" gorm:"-"`
		} `json:"accountVerified"`
		ScreenNameChanges []struct {
			ScreenNameChange struct {
				AccountID        int              `json:"accountId,string"`
				ScreenNameChange ScreenNameChange `json:"screenNameChange"`
			} `json:"screenNameChange"`
		}
		EmailAddressChanges []struct {
			EmailAddressChange struct {
				AccountID          int                `json:"accountId,string"`
				EmailAddressChange EmailAddressChange `json:"emailChange"`
			} `json:"emailAddressChange"`
		}
		IPAudit []struct {
			LoginIP LoginIP `json:"ipAudit"`
		}
	}, 1)

	if err := readJSON(
		path.Join(inputPath, "account.js"),
		"window.YTD.account.part0 = ",
		&data,
	); err != nil {
		return err
	}

	if err := readJSON(
		path.Join(inputPath, "account-creation-ip.js"),
		"window.YTD.account_creation_ip.part0 = ",
		&data,
	); err != nil {
		return err
	}

	if err := readJSON(
		path.Join(inputPath, "account-timezone.js"),
		"window.YTD.account_timezone.part0 = ",
		&data,
	); err != nil {
		return err
	}

	if err := readJSON(
		path.Join(inputPath, "verified.js"),
		"window.YTD.verified.part0 = ",
		&data,
	); err != nil {
		return err
	}

	if err := readJSON(
		path.Join(inputPath, "screen-name-change.js"),
		"window.YTD.screen_name_change.part0 = ",
		&data[0].ScreenNameChanges,
	); err != nil {
		return err
	}

	if err := readJSON(
		path.Join(inputPath, "email-address-change.js"),
		"window.YTD.email_address_change.part0 = ",
		&data[0].EmailAddressChanges,
	); err != nil {
		return err
	}

	if err := readJSON(
		path.Join(inputPath, "ip-audit.js"),
		"window.YTD.ip_audit.part0 = ",
		&data[0].IPAudit,
	); err != nil {
		return err
	}

	account := data[0].Account

	if data[0].AccountCreationIP.AccountId == account.ID {
		account.accountCreationIP = data[0].AccountCreationIP.accountCreationIP
	}

	if data[0].AccountTimezone.AccountId == account.ID {
		account.accountTimezone = data[0].AccountTimezone.accountTimezone
	}

	if data[0].AccountVerified.AccountId == account.ID {
		account.accountVerified = data[0].AccountVerified.accountVerified
	}

	for _, entity := range data[0].ScreenNameChanges {
		screenNameChange := entity.ScreenNameChange.ScreenNameChange
		screenNameChange.AccountID = entity.ScreenNameChange.AccountID

		if screenNameChange.AccountID == account.ID {
			err := p.DB().
				FirstOrCreate(&screenNameChange, screenNameChange.Conditions()).
				Error
			if err != nil {
				return err
			}

			account.ScreenNameChanges = append(account.ScreenNameChanges, screenNameChange)
		}
	}

	for _, entity := range data[0].EmailAddressChanges {
		emailAddressChange := entity.EmailAddressChange.EmailAddressChange
		emailAddressChange.AccountID = entity.EmailAddressChange.AccountID

		if emailAddressChange.AccountID == account.ID {
			err := p.DB().
				FirstOrCreate(&emailAddressChange, emailAddressChange.Conditions()).
				Error
			if err != nil {
				return err
			}

			account.EmailAddressChanges = append(account.EmailAddressChanges, emailAddressChange)
		}
	}

	for _, entity := range data[0].IPAudit {
		loginIP := entity.LoginIP

		if loginIP.AccountID == account.ID {
			err := p.DB().
				FirstOrCreate(&loginIP, loginIP.Conditions()).
				Error
			if err != nil {
				return err
			}

			account.LoginIPs = append(account.LoginIPs, loginIP)
		}
	}

	err := p.DB().
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&account).
		Error
	if err != nil {
		return err
	}

	return nil
}
