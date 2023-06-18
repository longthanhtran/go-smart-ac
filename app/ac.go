package app

import "time"

type Ac struct {
	Serial          string    `json:"serial"`
	RegisterDate    time.Time `json:"register_date"`
	FirmwareVersion string    `json:"firmware_version"`
	Token           string    `json:"token"`
}
