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

// RegenerateIndex rebuilds main index with saved directory
func RegenerateIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	index.I.Regenerate()
	log.WInfo(w, "regenerated index")
}

// DeleteKey deletes the file associated with the given key, returns 404 if not found
func DeleteKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName("key")
	log.Info("delete key '%s'", key)
	file, ok := index.I.Lookup(key)

	// if file found delete it
	if ok {
		err := index.I.Delete(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.WWarn(w, "err unable to delete key '%s': '%s'", key, err.Error())
			return
		}

		log.WInfo(w, "delete '%s' successful", key)
		return
	}

	// else state not found
	w.WriteHeader(http.StatusNotFound)
	log.WWarn(w, "key '%s' does not exist", key)
}
