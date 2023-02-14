package index

import (
	"fmt"
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

// Lookup returns the file with that key
// Returns (File, true) if file exists
// otherwise, returns new File, false
func (i *FileIndex) Lookup(key string) (*File, bool) {
	// read lock on index
	i.mu.RLock()
	defer i.mu.RUnlock()

	// get if File exists, return nil and false otherwise
	if file, ok := i.index[key]; ok {
		return file, true
	}

	return &File{FileName: key}, false
}

// Put creates/updates file in the fileindex
func (i *FileIndex) Put(file *File, bytes []byte) error {
	// write lock on index
	i.mu.Lock()
	defer i.mu.Unlock()

	i.index[file.FileName] = file
	err := file.ReplaceContent(string(bytes))
	return err
}

// ResolvePath returns a string representing the path to file
func (f *File) ResolvePath() string {
	if I.dir == "" {
		return fmt.Sprintf("%s.json", f.FileName)
	}

	return fmt.Sprintf("%s/%s.json", I.dir, f.FileName)
}

// Delete deletes the given file and then removes it from I
func (i *FileIndex) Delete(file *File) error {
	// write lock on index
	i.mu.Lock()
	defer i.mu.Unlock()

	// delete first so pointer isn't nil
	err := file.Delete()
	if err != nil {
		delete(i.index, file.FileName)
	}

	return err
}
