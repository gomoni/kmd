package ocr

import (
	"fmt"
	"io"

	ocr "github.com/otiai10/gosseract/v2"
)

type Client struct {
	client *ocr.Client
}

func NewClient() Client {
	c := Client{
		client: ocr.NewClient(),
	}
	c.Languages([]string{"eng"})
	return c
}

func (c Client) Close() {
	c.client.Close()
}

func (c Client) ImagePath(path string) error {
	return c.client.SetImage(path)
}

func (c Client) ImageBytes(data []byte) error {
	return c.client.SetImageFromBytes(data)
}

func (c Client) ImageReader(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}
	return c.client.SetImageFromBytes(data)
}

func (c Client) Languages(languages []string) {
	c.client.Languages = languages
}

func (c Client) Text() (string, error) {
	return c.client.Text()
}
