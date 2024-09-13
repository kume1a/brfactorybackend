package main

import (
	"brfactorybackend/internal/config"
	"brfactorybackend/internal/shared"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/robfig/cron"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		config.CreateCollections(app)

		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./public"), false))

		e.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
			AllowHeaders: []string{"Content-Type", "Authorization"},
		}))

		scheduler := cron.New()

		scheduler.AddFunc("*/10 * * * *", func() {
			log.Println("Sending email")
			err := shared.SendEmail(app, shared.SendEmailArgs{
				ToEmail: "kumela011@gmail.com",
				Subject: "Alert",
				Text:    "Hello",
			})

			if err != nil {
				log.Println("Error sending an email, ", err)
			}
		})

		scheduler.Start()

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
