package main

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/longthanhtran/go-smart-ac/api"
	"github.com/longthanhtran/go-smart-ac/database"
	"regexp"
	"strings"

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
	err = database.DBConn.AutoMigrate(&database.Ac{})
	err = database.DBConn.AutoMigrate(&database.Status{})
	if err != nil {
		return
	}
	fmt.Println("Database migrated")
}

var (
	protectedUrls = []*regexp.Regexp{
		regexp.MustCompile("/api/acs/[a-z0-9]*"),
	}
)

func validateAPIKey(_ *fiber.Ctx, key string) (bool, error) {
	db := database.DBConn
	var ac database.Ac
	result := db.First(&ac, "token = ?", key)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, keyauth.ErrMissingOrMalformedAPIKey
	}
	return true, nil
}

func authFilter(c *fiber.Ctx) bool {
	originalURL := strings.ToLower(c.OriginalURL())

	for _, pattern := range protectedUrls {
		if pattern.MatchString(originalURL) {
			return false
		}
	}
	return true
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(keyauth.New(keyauth.Config{
		Next:      authFilter,
		Validator: validateAPIKey,
	}))

	initDb()

	acApi := app.Group("/api/acs")
	acApi.Post("/", api.Create)
	acApi.Get("/", api.Show)
	acApi.Post("/status", api.StatusUpdate)
	acApi.Post("/status/bulk", api.BulkUpdate)

	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
