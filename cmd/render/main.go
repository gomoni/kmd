package main

import (
	"image/png"
	"log"
	"os"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
)

// Be sure to close pools/instances when you're done with them.
var pool pdfium.Pool

func init() {
	// Init the PDFium library and return the instance to open documents.
	// You can tweak these configs to your need. Be aware that workers can use quite some memory.
	var err error
	pool, err = webassembly.Init(webassembly.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	const input = `testdata/test.pdf`
	err := maine(input, "test.png")
	if err != nil {
		log.Fatal(err)
	}
}

func maine(input, output string) error {
	instance, err := pool.GetInstance(time.Second * 30)
	if err != nil {
		return err
	}
	defer instance.Close()

	doc, err := instance.OpenDocument(&requests.OpenDocument{FilePath: &input})
	if err != nil {
		return err
	}
	// Always close the document, this will release its resources.
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	// Render the page in DPI 200.
	pageRender, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
		DPI: 200, // The DPI to render the page in.
		Page: requests.Page{
			ByIndex: &requests.PageByIndex{
				Document: doc.Document,
				Index:    0,
			},
		}, // The page to render, 0-indexed.
	})
	if err != nil {
		return err
	}

	// The Render* methods return a cleanup function that has to be called when
	// using webassembly to make sure resources are cleaned up. Do this after
	// you are done with the returned image object.
	defer pageRender.Cleanup()

	// Write the output to a file.
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, pageRender.Result.Image)
	if err != nil {
		return err
	}

	return nil
}
