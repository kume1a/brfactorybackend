package main

import (
	"brfactorybackend/internal/config"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
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

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
