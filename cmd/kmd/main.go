// A trivial server for a tesseract-ocr
package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/coreos/go-systemd/activation"
	kos "github.com/gomoni/kmd/internal/os"
	"github.com/gomoni/kmd/internal/render"
	"github.com/gomoni/kmd/internal/server"
)

const (
	maxMemory = 32 << 20 // a limit on uploaded file size
)

type Globals struct {
	host string
}

type CLI struct {
	Host  string   `optional:"" short:"H" env:"KMD_HOST" help:"Server host, defaults to $XDG_RUNTIME_DIR/kmd.sock or to KMD_HOST environment variable."`
	Serve ServeCmd `cmd:"" default:"1" help:"run the server"`
}

type ServeCmd struct{}

func (cmd *ServeCmd) Run(g Globals) error {
	pool, err := render.NewPool()
	if err != nil {
		log.Fatalf("NewPool err: %s", err)
	}
	defer pool.Close()

	ocr := server.NewOCR(maxMemory, pool)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", server.Info)
	mux.Handle("POST /ocr", ocr)

	var listener net.Listener
	if _, socketActivated := os.LookupEnv("LISTEN_FDS"); socketActivated {
		listeners, err := activation.Listeners()
		if err != nil {
			log.Fatal(err)
		}
		if listeners == nil {
			log.Fatal("no listeners passed, yet LISTEN_FDS defined")
		}
		listener = listeners[0]
	} else {
		l, err := net.Listen("unix", g.host)
		if err != nil {
			log.Fatal(err)
		}
		listener = l
		defer os.Remove(g.host)
	}

	s := &http.Server{
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s.Serve(listener)
}

func maine() error {
	var cli CLI
	parsed := kong.Parse(&cli)
	if cli.Host == "" {
		cli.Host = kos.DefaultSocketPath()
	}
	return parsed.Run(&cli, Globals{
		host: cli.Host,
	})
}

func main() {
	if err := maine(); err != nil {
		log.Fatal(err)
	}
}
