package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/longthanhtran/go-smart-ac/database"
	"gorm.io/gorm"
	"strings"
)

func Create(c *fiber.Ctx) error {
	db := database.DBConn
	ac := new(database.Ac)
	if err := c.BodyParser(ac); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	ac.Token = createToken()
	db.Create(&ac)
	return c.JSON(ac.ToJson())
}

func Show(c *fiber.Ctx) error {
	db := database.DBConn
	_, serial := acData(c, db)
	var ac database.Ac
	db.First(&ac, "serial = ?", serial)
	return c.JSON(ac.ToJson())
}

func StatusUpdate(c *fiber.Ctx) error {
	db := database.DBConn

	acStatus := new(database.Status)
	if err := c.BodyParser(acStatus); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	valid, acSerial := acData(c, db)
	if acSerial != acStatus.AcSerial || !valid {
		return c.Status(401).SendString("unauthorized")
	}
	if !validateStatus(*acStatus) {
		return c.Status(400).SendString("invalid ac status")
	}
	db.Create(&acStatus)
	return c.JSON(acStatus.ToJson())
}

func acData(c *fiber.Ctx, db *gorm.DB) (bool, string) {
	token := strings.Split(c.GetReqHeaders()["Authorization"], " ")[1]
	var ac database.Ac
	result := db.First(&ac, "token = ?", token)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, ""
	}
	return true, ac.Serial
}

func BulkUpdate(c *fiber.Ctx) error {
	db := database.DBConn

	// unmarshal []status from c.Body()
	var statuses []database.Status
	acs := c.Body()
	if err := json.Unmarshal(acs, &statuses); err != nil {
		return c.SendStatus(400)
	}

	if len(statuses) > 500 {
		return c.Status(400).SendString("more than 500 statuses sent")
	}
	_, serial := acData(c, db)
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

func createToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
