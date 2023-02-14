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

// NewFileIndex returns a reference to a new file index
func NewFileIndex(dir string) *FileIndex {
	return &FileIndex{
		dir:        dir,
		index:      map[string]*File{},
		FileSystem: af.NewOsFs(),
	}
}

// SetFileSystem sets the file system for the given FileIndex
func (i *FileIndex) SetFileSystem(fs af.Fs) {
	i.FileSystem = fs
}

// List all keys in database
func (i *FileIndex) List() (res []string) {
	// read lock on index
	i.mu.RLock()
	defer i.mu.RUnlock()

	for k := range i.index {
		res = append(res, k)
	}

	return res
}
