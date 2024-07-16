package render

import (
	"fmt"
	"image/png"
	"io"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
)

type Pool struct {
	pool pdfium.Pool
}

func NewPool() (Pool, error) {
	pool, err := webassembly.Init(webassembly.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
	})
	if err != nil {
		return Pool{}, fmt.Errorf("initializing web assembly pool: %w", err)
	}
	return Pool{pool: pool}, nil
}

func (p Pool) Close() error {
	return p.pool.Close()
}

// Render renders first page from the pdf document as a png stream
func (p Pool) Render(tmout time.Duration, r io.ReadSeeker, size int64, w io.Writer) (err error) {
	instance, err := p.pool.GetInstance(tmout)
	if err != nil {
		return fmt.Errorf("getting instance from pool: %w", err)
	}
	defer instance.Close()

	doc, err := instance.OpenDocument(&requests.OpenDocument{
		FileReader:     r,
		FileReaderSize: size,
	})
	if err != nil {
		return fmt.Errorf("opening document: %w", err)
	}
	// Always close the document, this will release its resources.
	defer func() {
		_, err = instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
			Document: doc.Document,
		})
	}()

	pageRender, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
		DPI: 200,
		Page: requests.Page{
			ByIndex: &requests.PageByIndex{
				Document: doc.Document,
				Index:    0,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("rendering page in 200 DPI: %w", err)
	}

	// The Render* methods return a cleanup function that has to be called when
	// using webassembly to make sure resources are cleaned up. Do this after
	// you are done with the returned image object.
	defer pageRender.Cleanup()

	err = png.Encode(w, pageRender.Result.Image)
	if err != nil {
		return fmt.Errorf("encoding png: %w", err)
	}

	return nil
}
