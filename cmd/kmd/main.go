// A trivial server for a tesseract-ocr
package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-systemd/activation"
	kos "github.com/gomoni/kmd/internal/os"
	"github.com/gomoni/kmd/internal/render"
	"github.com/gomoni/kmd/internal/server"
)

const (
	maxMemory = 32 << 20 // a limit on uploaded file size
)

func main() {
	pool, err := render.NewPool()
	if err != nil {
		log.Fatalf("NewPool err: %s", err)
	}
	defer pool.Close()

	sock := kos.DefaultSocketPath()

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
		l, err := net.Listen("unix", sock)
		if err != nil {
			log.Fatal(err)
		}
		listener = l
		defer os.Remove(sock)
	}

	s := &http.Server{
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.Serve(listener); err != nil {
		panic(err)
	}
}
