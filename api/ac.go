package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/longthanhtran/go-smart-ac/database"
	"gorm.io/gorm"
)

func Create(c *fiber.Ctx) error {
	db := database.DBConn
	ac := new(database.Ac)
	if err := c.BodyParser(ac); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	db.Create(&ac)
	return c.JSON(ac)
}

func Show(c *fiber.Ctx) error {
	db := database.DBConn
	acSerial := c.Params("serial")
	var ac database.Ac
	db.First(&ac, "serial = ?", acSerial)
	return c.JSON(ac.AcToRespondAc())
}

func StatusUpdate(c *fiber.Ctx) error {
	db := database.DBConn
	acSerial := c.Params("serial")
	var ac database.Ac
	foundAc := db.First(&ac, "serial = ?", acSerial)
	if foundAc.Error != nil && errors.Is(foundAc.Error, gorm.ErrRecordNotFound) {
		return c.SendStatus(404)
	}
	acStatus := new(database.Status)
	if err := c.BodyParser(acStatus); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	db.Create(&acStatus)
	return c.JSON(acStatus.StatusToAcStatus())
}
