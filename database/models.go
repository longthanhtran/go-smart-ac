package database

import (
	"github.com/longthanhtran/go-smart-ac/app"
	"gorm.io/gorm"
	"time"
)

type Ac struct {
	gorm.Model
	Serial          string    `json:"serial" gorm:"uniqueIndex"`
	RegisterDate    time.Time `json:"register_date"`
	FirmwareVersion string    `json:"firmware_version"`
}

func (ac *Ac) ToJson() app.Ac {
	return app.Ac{
		Serial:          ac.Serial,
		RegisterDate:    ac.RegisterDate,
		FirmwareVersion: ac.FirmwareVersion,
	}
}

type Status struct {
	gorm.Model
	AcSerial     string `json:"ac_serial"`
	Serial       Ac     `gorm:"foreignKey:AcSerial"`
	Temperature  uint8  `json:"temperature"`
	Humidity     uint8  `json:"humidity"`
	CoLevel      uint16 `json:"co_level"`
	HealthStatus string `json:"health_status" gorm:"size:150" validate:"max=150"`
}

func (status *Status) ToJson() app.Status {
	return app.Status{
		AcSerial:     status.AcSerial,
		Temperature:  status.Temperature,
		Humidity:     status.Humidity,
		CoLevel:      status.CoLevel,
		HealthStatus: status.HealthStatus,
	}
}
