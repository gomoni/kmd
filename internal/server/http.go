package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gomoni/kmd/internal/ocr"
)

type PDFRenderer interface {
	Render(tmout time.Duration, r io.ReadSeeker, size int64, w io.Writer) (err error)
}

type OCR struct {
	maxMemory int64
	renderer  PDFRenderer
}

func NewOCR(maxMemory int64, renderer PDFRenderer) OCR {
	return OCR{
		maxMemory: maxMemory,
		renderer:  renderer,
	}
}

func (o OCR) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(o.maxMemory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	upload, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer upload.Close()

	client := ocr.NewClient()
	defer client.Close()

	sc := ocr.NewSmartClient(client, o.renderer)

	err = sc.ImageReader(upload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// argument handling
	if langs := r.FormValue("languages"); langs != "" {
		client.Languages(strings.Split(langs, ","))
	}

	out, err := client.Text()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, out)
}

func Info(w http.ResponseWriter, r *http.Request) {
	info, err := ocr.Info()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: application/json response too
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "version: %s\n", info.Version)
	fmt.Fprintf(w, "languages:\n")
	for _, l := range info.Languages {
		fmt.Fprintf(w, " * %s\n", l)
	}
}
