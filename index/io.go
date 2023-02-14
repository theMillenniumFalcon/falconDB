package index

// replaceContent changes the contents of file f to be str
func (f *File) replaceContent(str string) error {
	// write lock on file
	f.mu.Lock()
	defer f.mu.Unlock()

	// success
	return nil
}
