package hfs

import "os"

// Sorter sorts file names for directory listing.
type Sorter interface {
	// Less returns a function for sort.Slice.
	Less([]os.FileInfo) func(int, int) bool
}
