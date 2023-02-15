// Package api contains files responsible for defining and serving the restful api
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

// GetKey returns the file with that key if found, otherwise return 404
func GetKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName("key")
	log.Info("get key '%s'", key)

	file, ok := index.I.Lookup(key)

	// if file fetch is successful
	if ok {
		w.Header().Set("Content-Type", "application/json")

		// unpack bytes into map
		jsonMap, err := file.ToMap()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.WWarn(w, "err key '%s' cannot be parsed into json: %s", key, err.Error())
			return
		}

		// successful field get
		w.Header().Set("Content-Type", "application/json")
		maxDepth := getMaxDepthParam(r)
		resolvedJsonMap := index.ResolveReferences(jsonMap, maxDepth)

		jsonData, _ := json.Marshal(resolvedJsonMap)
		fmt.Fprintf(w, "%+v", string(jsonData))
		return
	}

	// otherwise write 404
	w.WriteHeader(http.StatusNotFound)
	log.WWarn(w, "key '%s' not found", key)
}

// try to find recursive depth param or else return a default
func getMaxDepthParam(r *http.Request) int {
	maxDepth := 3

	maxDepthStr := r.URL.Query().Get("depth")
	parsedInt, err := strconv.Atoi(maxDepthStr)
	if err == nil {
		maxDepth = parsedInt
	}

	return maxDepth
}
