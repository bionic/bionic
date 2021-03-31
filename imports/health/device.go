package health

import (
	"gorm.io/gorm"
	"strings"
)

type Device struct {
	gorm.Model
	Name         string
	Manufacturer string
	DeviceModel  string `gorm:"column:model"`
	Hardware     string
	Software     string
}

func (Device) TableName() string {
	return tablePrefix + "devices"
}

func (d Device) Conditions() map[string]interface{} {
	return map[string]interface{}{
		"name":         d.Name,
		"manufacturer": d.Manufacturer,
		"model":        d.DeviceModel,
		"hardware":     d.Hardware,
		"software":     d.Software,
	}
}

func (d *Device) UnmarshalText(b []byte) error {
	text := string(b)

	if len(text) < 3 {
		return nil
	}

	parts := strings.Split(text[1:len(text)-1], ", ")
	if len(parts) < 2 {
		return nil
	}

	attributes := parts[1:]

	for _, attr := range attributes {
		attrParts := strings.Split(attr, ":")
		if len(attrParts) != 2 {
			continue
		}

		key, value := attrParts[0], attrParts[1]

		switch key {
		case "name":
			d.Name = value
		case "manufacturer":
			d.Manufacturer = value
		case "model":
			d.DeviceModel = value
		case "hardware":
			d.Hardware = value
		case "software":
			d.Software = value
		}
	}

	return nil
}
