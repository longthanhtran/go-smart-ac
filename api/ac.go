package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/longthanhtran/go-smart-ac/database"
	"gorm.io/gorm"
	"time"
)

type Ac struct {
	gorm.Model
	Serial          string    `json:"serial" gorm:"uniqueIndex"`
	RegisterDate    time.Time `json:"register_date"`
	FirmwareVersion string    `json:"firmware_version"`
}

type Status struct {
	gorm.Model
	AcSerial     string `json:"ac_serial"`
	Serial       Ac     `gorm:"foreignKey:AcSerial"`
	Temperature  uint8  `json:"temperature"`
	Humidity     uint8  `json:"humidity"`
	CoLevel      uint16 `json:"co_level"`
	HealthStatus string `json:"health_status" gorm:"size:150"`
}

func Create(c *fiber.Ctx) error {
	db := database.DBConn
	ac := new(Ac)
	if err := c.BodyParser(ac); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	db.Create(&ac)
	return c.JSON(ac)
}

func StatusUpdate(c *fiber.Ctx) error {
	db := database.DBConn
	acId := c.Params("ac_id")
	var ac Ac
	foundAc := db.First(&ac, "serial = ?", acId)
	if foundAc.Error != nil && errors.Is(foundAc.Error, gorm.ErrRecordNotFound) {
		return c.SendStatus(404)
	}
	acStatus := new(Status)
	if err := c.BodyParser(acStatus); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	db.Create(&acStatus)
	return c.JSON(acStatus)
}
