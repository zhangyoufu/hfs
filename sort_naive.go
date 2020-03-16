package hfs

import (
	"os"
	"strings"
)

// NaiveSorter sorts file names using Unicode codepoints.
type NaiveSorter struct {
	// DirectoryFirst put directories before files.
	DirectoryFirst bool

	// IgnoreCase makes use of strings.ToUpper before comparing Unicode
	// codepoints.
	IgnoreCase bool
}

// Less returns a function for sort.Slice.
func (s NaiveSorter) Less(files []os.FileInfo) func(int, int) bool {
	return func(i, j int) bool {
		if s.DirectoryFirst {
			i_dir := files[i].IsDir()
			j_dir := files[j].IsDir()
			if i_dir && !j_dir {
				return true
			}
			if !i_dir && j_dir {
				return false
			}
		}
		i_name := files[i].Name()
		j_name := files[j].Name()
		if s.IgnoreCase {
			i_name = strings.ToUpper(i_name)
			j_name = strings.ToUpper(j_name)
		}
		return i_name < j_name
	}
}
