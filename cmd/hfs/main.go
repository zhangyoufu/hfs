package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	proxyproto "github.com/pires/go-proxyproto"
	"github.com/zhangyoufu/hfs"
)

func main() {
	var (
		network     string
		address     string
		root        string
		enableProxy bool
		accessLog   bool
		dirList     bool
		dirFirst    bool
		ignoreCase  bool
		dotFile     bool
		indexPage   bool
		protoH1     bool
		protoH2C    bool
	)

	flag.StringVar(&network, "network", "tcp", "tcp or unix")
	flag.StringVar(&address, "addr", "localhost:8000", "address:port for tcp, or filesystem path for unix")
	flag.StringVar(&root, "root", ".", "`path` of document root")
	flag.BoolVar(&enableProxy, "proxy", false, "enable PROXY protocol support (default false)")
	flag.BoolVar(&accessLog, "log", true, "enable access log")
	flag.BoolVar(&dirList, "list", true, "enable directory listing")
	flag.BoolVar(&dirFirst, "dirfirst", true, "list directories before files")
	flag.BoolVar(&ignoreCase, "igncase", true, "case insensitive sorting")
	flag.BoolVar(&dotFile, "dotfile", false, "enable listing and serving dot files (default false)")
	flag.BoolVar(&indexPage, "index", true, "enable serving index.html")
	flag.BoolVar(&protoH1, "h1", true, "enable HTTP/1 protocol")
	flag.BoolVar(&protoH2C, "h2c", false, "enable unencrypted HTTP/2 protocol (default false)")
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

	switch network {
	case "tcp":
	case "unix":
		_ = os.Remove(address)
		_ = umask(0o007)
	default:
		log.Fatal("unsupported network")
	}

	ln, err := net.Listen(network, address)
	if err != nil {
		log.Fatal("unable to listen: ", err)
	}

	if enableProxy {
		ln = &proxyproto.Listener{
			Listener:   ln,
			ConnPolicy: func (proxyproto.ConnPolicyOptions) (proxyproto.Policy, error) {
				return proxyproto.REQUIRE, nil
			},
			// ReadHeaderTimeout defaults to 10s
		}
	}

	protocols := http.Protocols{}
	protocols.SetHTTP1(protoH1)
	protocols.SetUnencryptedHTTP2(protoH2C)
	if protocols.String() == "{}" {
		log.Fatal("no protocols available")
	}

	server := http.Server{
		Handler: &hfs.FileServer{
			FileSystem:       http.Dir(root),
			Sorter:           sorter,
			AccessLog:        logger,
			ErrorLog:         log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds),
			DirectoryListing: dirList,
			ServeDotFile:     dotFile,
			ServeIndexPage:   indexPage,
		},
		DisableGeneralOptionsHandler: true,
		Protocols: &protocols,
	}

	log.Printf("Serving %s on [%s] %s %s", root, network, address, protocols)
	log.Fatal(server.Serve(ln))
}
