package main

import (
	"brfactorybackend/internal/modules/scheduledwork"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "brfactorybackend/migrations"
)

func main1() {
	app := pocketbase.New()

	// err := config.LoadEnv(app)
	// if err != nil {
	// 	log.Fatal("Couldn't load env vars, returning")
	// }

	isGoRun := strings.Contains(os.Args[0], "\\tmp\\") || strings.Contains(os.Args[0], "/tmp/")
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./public"), false))

		e.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
			AllowHeaders: []string{"Content-Type", "Authorization"},
		}))

		scheduledwork.SetupCronJobs(app)

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
