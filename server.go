package hfs

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
)

// A modified version of net/http.FileServer.
type FileServer struct {
	// Abstraction of filesystem
	FileSystem http.FileSystem

	// Sorter can be nil if DirectoryListing is disabled
	Sorter Sorter

	// Logger for access log
	AccessLog *log.Logger

	// Logger for error log
	ErrorLog *log.Logger

	// Enable directory listing
	DirectoryListing bool

	// Allow listing and serving dotfiles
	ServeDotFile bool

	// Allow serving index.html
	ServeIndexPage bool
}

func (s *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logf(s.AccessLog, "%s %s", stripPort(r.RemoteAddr), r.URL.EscapedPath())

	// check path
	urlPath := r.URL.Path
	valid, endWithSlash := true, false
	switch {
	case len(urlPath) < 1:
		valid = false
	case urlPath[0] != '/':
		valid = false
	default:
		endWithSlash = urlPath[len(urlPath)-1] == '/'
		if endWithSlash {
			urlPath = urlPath[:len(urlPath)-1]
		}
		if urlPath == "" {
			break
		}
		for _, elem := range strings.Split(urlPath[1:], "/") {
			if elem == "" || elem == "." || elem == ".." || !s.ServeDotFile && elem[0] == '.' {
				valid = false
				break
			}
		}
	}
	if !valid {
		badRequest(w)
		return
	}

	// open
	file, err := s.FileSystem.Open(urlPath) // path will be checked by internal/safefilepath
	if err != nil {
		switch {
		case os.IsNotExist(err):
			notFound(w)
		case os.IsPermission(err):
			logf(s.ErrorLog, "failed to open %s: %v", urlPath, err)
			forbidden(w)
		case err.Error() == "http: invalid or unsafe file path":
			// all platform: NUL
			// Windows-only: backslash(\), colon(:)
			logf(s.ErrorLog, "bad path")
			badRequest(w)
		default:
			logf(s.ErrorLog, "failed to open %s: %v", urlPath, err)
			internalError(w)
		}
		return
	}
	defer file.Close()

	// stat
	stat, err := file.Stat()
	if err != nil {
		logf(s.ErrorLog, "failed to stat %s: %v", urlPath, err)
		internalError(w)
		return
	}

	if endWithSlash {
		// serve directory

		if !stat.IsDir() {
			notFound(w)
			return
		}

		if s.ServeIndexPage {
			const indexPage = "/index.html"
			indexPath := urlPath + indexPage
			indexFile, err := s.FileSystem.Open(indexPath)
			switch {
			case err == nil:
				// serve index page
				defer indexFile.Close()
				indexStat, err := indexFile.Stat()
				if err != nil {
					logf(s.ErrorLog, "failed to stat %s: %v", indexPath, err)
					internalError(w)
					return
				}
				if indexStat.IsDir() {
					notFound(w)
					return
				}
				http.ServeContent(w, r, indexPage, indexStat.ModTime(), indexFile)
				return
			case os.IsPermission(err):
				logf(s.ErrorLog, "failed to open %s: %v", indexPath, err)
				forbidden(w)
				return
			case os.IsNotExist(err):
				break
			default:
				logf(s.ErrorLog, "failed to open %s: %v", indexPath, err)
				internalError(w)
				return
			}
			// fall through
		}

		if s.DirectoryListing {
			s.dirList(w, file, urlPath)
			return
		}

		notFound(w)
		return
	} else {
		if stat.IsDir() {
			if s.DirectoryListing {
				redirect(w, path.Base(urlPath)+"/")
				return
			} else {
				notFound(w)
				return
			}
		}

		// serve file
		http.ServeContent(w, r, urlPath, stat.ModTime(), file)
		return
	}
}

func (s *FileServer) dirList(w http.ResponseWriter, f http.File, path string) {
	top := false
	if path == "" {
		path = "/"
		top = true
	}

	files, err := f.Readdir(-1)
	if err != nil {
		logf(s.ErrorLog, "failed to list %s: %v", path, err)
		internalError(w)
		return
	}

	i := 0
	for _, file := range files {
		name := file.Name()
		if name == "" {
			continue
		} // sanity check
		if !s.ServeDotFile && name[0] == '.' {
			continue
		}
		files[i] = file
		i++
	}
	files = files[:i]

	sort.Slice(files, s.Sorter.Less(files))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")
	if !top {
		fmt.Fprintln(w, `<a href="../">../</a>`)
	}
	for _, entry := range files {
		name := entry.Name()
		escapedName := url.PathEscape(name)
		if entry.IsDir() {
			name += "/"
			escapedName += "/"
		}
		fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", escapedName, html.EscapeString(name))
	}
	fmt.Fprintf(w, "</pre>\n")
}
