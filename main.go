package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/longthanhtran/go-smart-ac/api"
	"github.com/longthanhtran/go-smart-ac/database"

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
	if err != nil {
		return
	}
	fmt.Println("Database migrated")
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	initDb()

	app.Post("/api/acs", api.Create)

	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
