// Package main serves as a user interface to the underlying index and api packages
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/theMillenniumFalcon/falconDB/api"
	"github.com/theMillenniumFalcon/falconDB/log"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "falconDB",
		Usage: "a in-memory document based database for fast prototyping",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dir",
				Aliases:     []string{"d"},
				Value:       "db",
				Usage:       "directory to look for keys",
				DefaultText: "db",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// serve defines all the endpoints and starts a new http server on :3000
func serve(port int, dir string) error {
	log.SetLoggingLevel(log.INFO)
	log.Info("initializing nanoDB")
	setup(dir)

	router := httprouter.New()

	// define endpoints
	router.GET("/", api.GetIndex)
	router.POST("/", api.RegenerateIndex)
	router.GET("/:key", api.GetKey)
	router.DELETE("/:key", api.DeleteKey)
	router.PATCH("/:key/:field", api.PatchKeyField)

	// start server
	log.Info("starting api server on port %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func setup(sir string) {}
