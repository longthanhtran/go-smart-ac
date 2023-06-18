package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
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
	return c.JSON(ac.ToJson())
}

func Show(c *fiber.Ctx) error {
	db := database.DBConn
	acSerial := c.Params("serial")
	var ac database.Ac
	db.First(&ac, "serial = ?", acSerial)
	return c.JSON(ac.ToJson())
}

func StatusUpdate(c *fiber.Ctx) error {
	db := database.DBConn
	valid, serial := validateAcSerialParam(c, db)
	if !valid {
		return c.Status(404).SendString("invalid ac serial")
	}
	acStatus := new(database.Status)
	if err := c.BodyParser(acStatus); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	if acStatus.AcSerial != serial || !validateStatus(*acStatus) {
		return c.Status(400).SendString("invalid ac status")
	}
	db.Create(&acStatus)
	return c.JSON(acStatus.ToJson())
}

func validateAcSerialParam(c *fiber.Ctx, db *gorm.DB) (bool, string) {
	acSerial := c.Params("serial")
	var ac database.Ac
	foundAc := db.First(&ac, "serial = ?", acSerial)
	if foundAc.Error != nil && errors.Is(foundAc.Error, gorm.ErrRecordNotFound) {
		return false, ""
	}
	return true, acSerial
}

func BulkUpdate(c *fiber.Ctx) error {
	db := database.DBConn
	serial := c.Params("serial")
	var ac database.Ac
	foundAc := db.First(&ac, "serial = ?", serial)
	if foundAc.Error != nil && errors.Is(foundAc.Error, gorm.ErrRecordNotFound) {
		return c.Status(404).SendString("invalid ac serial")
	}

	// unmarshal []status from c.Body()
	var statuses []database.Status
	acs := c.Body()
	if err := json.Unmarshal(acs, &statuses); err != nil {
		return c.SendStatus(400)
	}

	if len(statuses) > 500 {
		return c.Status(400).SendString("more than 500 statuses sent")
	}
	for _, status := range statuses {
		if status.AcSerial != serial || !createStatus(status, db) {
			fmt.Printf("invalid ac status of %s\n", status.AcSerial)
		}
	}
	return c.SendStatus(200)
}

func createStatus(status database.Status, db *gorm.DB) bool {
	if !validateStatus(status) {
		return false
	}
	db.Create(&status)
	return true
}

func validateStatus(status database.Status) bool {
	var validate = validator.New()

	if err := validate.Struct(status); err != nil {
		return false
	}
	return true
}
