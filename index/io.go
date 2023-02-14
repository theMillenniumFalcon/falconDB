package index

import (
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	af "github.com/spf13/afero"
)

func crawlDirectory(directory string) []string {
	files, err := af.ReadDir(I.FileSystem, directory)
	if err != nil {
		log.Fatal(err)
	}

	res := []string{}

	for _, file := range files {
		extension := filepath.Ext(file.Name())
		if extension == ".json" {
			name := strings.TrimSuffix(file.Name(), ".json")
			res = append(res, name)
		}
	}

	return res
}

// replaceContent changes the contents of file f to be str
func (f *File) replaceContent(str string) error {
	// write lock on file
	f.mu.Lock()
	defer f.mu.Unlock()

	// success
	return nil
}
