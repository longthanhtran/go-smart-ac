package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/longthanhtran/go-smart-ac/database"
	"gorm.io/gorm"
	"time"
)

type Ac struct {
	gorm.Model
	Serial          string    `json:"serial"`
	RegisterDate    time.Time `json:"register_date"`
	FirmwareVersion string    `json:"firmware_version"`
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
