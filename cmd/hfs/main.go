package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/zhangyoufu/hfs"
)

func main() {
	var (
		addr, root                                                   string
		accessLog, dirList, dirFirst, ignoreCase, dotFile, indexPage bool
	)

	flag.StringVar(&addr, "addr", "localhost:8000", "listen `address:port`")
	flag.StringVar(&root, "root", ".", "`path` of document root")
	flag.BoolVar(&accessLog, "log", true, "enable access log")
	flag.BoolVar(&dirList, "list", true, "enable directory listing")
	flag.BoolVar(&dirFirst, "dirfirst", true, "list directories before files")
	flag.BoolVar(&ignoreCase, "igncase", true, "case insensitive sorting")
	flag.BoolVar(&dotFile, "dotfile", false, "enable listing and serving dot files (default false)")
	flag.BoolVar(&indexPage, "index", true, "enable serving index.html")
	flag.Parse()

	sorter := hfs.Sorter(nil)
	if dirList {
		sorter = hfs.NaiveSorter{
			DirectoryFirst: dirFirst,
			IgnoreCase:     ignoreCase,
		}
	}

	var logger *log.Logger
	if accessLog {
		logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	}

	log.Printf("Serving %s on %s", root, addr)
	log.Fatal(http.ListenAndServe(addr, &hfs.FileServer{
		FileSystem:       http.Dir(root),
		Sorter:           sorter,
		AccessLog:        logger,
		ErrorLog:         log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds),
		DirectoryListing: dirList,
		ServeDotFile:     dotFile,
		ServeIndexPage:   indexPage,
	}))
}
