// Package main serves as a user interface to the underlying index and api packages
package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/theMillenniumFalcon/falconDB/log"
)

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
	router.GET("/:key/:field", api.GetKeyField)
	router.PUT("/:key", api.UpdateKey)
	router.DELETE("/:key", api.DeleteKey)
	router.PATCH("/:key/:field", api.PatchKeyField)

	// start server
	log.Info("starting api server on port %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
