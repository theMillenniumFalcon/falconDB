package index

import (
	"sync"

	af "github.com/spf13/afero"
)

// I is the global database index which keeps track of
// which files are where
var I *FileIndex

// FileIndex is holds the actual index mapping for keys to files
type FileIndex struct {
	mu         sync.RWMutex
	dir        string
	index      map[string]*File
	FileSystem af.Fs
}

// File stores the filename as well as a read-write mutex
type File struct {
	FileName string
	mu       sync.RWMutex
}
