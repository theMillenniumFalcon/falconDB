// Package api contains files responsible for defining and serving the restful api
package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/theMillenniumFalcon/falconDB/index"
	"github.com/theMillenniumFalcon/falconDB/log"
)

// GetIndex returns a JSON of all files in db index
func GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Info("retrieving index")
	files := index.I.List()

	// create temporary struct with index data
	data := struct {
		Files []string `json:"files"`
	}{
		Files: files,
	}

	// create json representation and return
	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(data)
	fmt.Fprintf(w, "%+v", string(jsonData))
}
