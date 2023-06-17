package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/longthanhtran/go-smart-ac/api"
	"github.com/longthanhtran/go-smart-ac/database"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDb() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("acs.sqlite"))
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connected to database")
	err = database.DBConn.AutoMigrate(&api.Ac{})
	err = database.DBConn.AutoMigrate(&api.Status{})
	if err != nil {
		return
	}
	fmt.Println("Database migrated")
}

var (
	apiKey = os.Getenv("X-API-KEY")
)

func validateAPIKey(_ *fiber.Ctx, key string) (bool, error) {
	hashedAPIKey := sha256.Sum256([]byte(apiKey))
	hashedKey := sha256.Sum256([]byte(key))

	if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
		return true, nil
	}
	return false, keyauth.ErrMissingOrMalformedAPIKey
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(keyauth.New(keyauth.Config{
		Validator: validateAPIKey,
	}))

	initDb()

	app.Post("/api/acs", api.Create)
	app.Post("/api/acs/:ac_id/status", api.StatusUpdate)

	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
