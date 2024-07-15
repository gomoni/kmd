package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gomoni/kmd/internal/ocr"
)

type OCR struct {
	maxMemory int64
}

func NewOCR(maxMemory int64) OCR {
	return OCR{
		maxMemory: maxMemory,
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

	err = client.ImageReader(upload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if langs := r.FormValue("languages"); langs != "" {
		client.Languages(strings.Split(langs, ","))
	}
	/*
		if whitelist := r.FormValue("whitelist"); whitelist != "" {
			err = client.SetWhitelist(whitelist)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	*/

	out, err := client.Text()
	if err != nil {
		// TODO: bad request error?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, out)
}
