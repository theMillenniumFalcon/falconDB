// Package main serves as a user interface to the underlying index and api packages
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
	"github.com/theMillenniumFalcon/falconDB/api"
	"github.com/theMillenniumFalcon/falconDB/index"
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
		Commands: []*cli.Command{
			{
				Name:    "start",
				Aliases: []string{"st"},
				Usage:   "start a falcondb server",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "port",
						Aliases:     []string{"p"},
						Value:       4000,
						Usage:       "port to run falcondb on",
						DefaultText: "4000",
					},
				},
				Action: func(c *cli.Context) error {
					return serve(c.Int("port"), c.String("dir"))
				},
			}, {
				Name:    "shell",
				Aliases: []string{"sh"},
				Usage:   "start an interactive falcondb shell",
				Action: func(c *cli.Context) error {
					return shell(c.String("dir"))
				},
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
	log.Info("initializing falconDB")
	setup(dir)

	router := httprouter.New()

	// define endpoints
	router.GET("/", api.GetIndex)
	router.POST("/", api.RegenerateIndex)
	router.GET("/:key", api.GetKey)
	router.GET("/:key/:field", api.GetKeyField)
	router.PUT("/:key", api.UpdateKey)
	router.DELETE("/:key", api.DeleteKey)
	router.PATCH("/:key/:field", api.PatchKeyField)

	// start server
	log.Info("starting api server on port %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func setup(dir string) {
	index.I = index.NewFileIndex(dir)

	// create falconDB lock
	err := makeLock(dir)
	if err != nil {
		log.Fatal(err)
		return
	}

	index.I.Regenerate()

	// trap sigint
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanUp(dir)
		os.Exit(1)
	}()
}

func cleanUp(dir string) {
	log.Info("\n caught term signal! cleaning up...")

	err := releaseLock(dir)
	if err != nil {
		log.Warn("couldn't remove lock")
		log.Fatal(err)
		return
	}
}

func makeLock(dir string) error {
	_, err := index.I.FileSystem.Stat(getLockLocation(dir))

	if os.IsNotExist(err) {
		_, err = index.I.FileSystem.Create(getLockLocation(dir))
		return err
	}

	return fmt.Errorf("couldn't acquire lock on %s", dir)
}

func releaseLock(dir string) error {
	lockdir := getLockLocation(dir)
	return index.I.FileSystem.Remove(lockdir)
}

func getLockLocation(dir string) string {
	base := "falconDB_lock"
	if dir == "" || dir == "." {
		return base
	}

	return dir + "/" + base
}
