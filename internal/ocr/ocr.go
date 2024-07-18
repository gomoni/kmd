package ocr

import (
	"bytes"
	"fmt"
	"io"
	"time"

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

func (c Client) ImageReader(r io.ReadSeeker) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}
	return c.client.SetImageFromBytes(data)
}

func (c Client) ImageBytes(data []byte) error {
	return c.client.SetImageFromBytes(data)
}

func (c Client) Languages(languages []string) {
	c.client.Languages = languages
}

func (c Client) Text() (string, error) {
	return c.client.Text()
}

type PDFRenderer interface {
	Render(tmout time.Duration, r io.ReadSeeker, size int64, w io.Writer) (err error)
}

type SmartClient struct {
	client   Client
	renderer PDFRenderer
}

func NewSmartClient(client Client, renderer PDFRenderer) SmartClient {
	return SmartClient{
		client:   client,
		renderer: renderer,
	}
}

func (sc SmartClient) ImageReader(r io.ReadSeeker) error {
	if isPDF, err := IsPDF(r); err != nil {
		return fmt.Errorf("detecting PDF: %w", err)
	} else if !isPDF {
		return sc.client.ImageReader(r)
	}

	size, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("SmartClient seek to end: %w", err)
	}
	png := bytes.NewBuffer(make([]byte, 0, size))
	err = sc.renderer.Render(time.Second*30, r, size, png)
	if err != nil {
		return fmt.Errorf("SmartClient render: %w", err)
	}
	err = sc.client.ImageBytes(png.Bytes())
	if err != nil {
		return fmt.Errorf("SmartClient ImageBytes: %w", err)
	}
	return nil
}

func (sc SmartClient) Text() (string, error) {
	ret, err := sc.client.Text()
	if err != nil {
		return "", fmt.Errorf("SmartClient Text: %w", err)
	}
	return ret, nil
}

// IsPDF checks if the reader is a PDF file via magic sequence %PDF-
// it DOES seek file back
func IsPDF(r io.ReadSeeker) (isPDF bool, err error) {
	seek0 := func() error {
		_, err = r.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("seek the input: %w", err)
		}
		return nil
	}

	if err = seek0(); err != nil {
		return false, err
	}
	defer func() {
		err = seek0()
	}()

	var magic [5]byte
	n, err := r.Read(magic[:])
	if err != nil {
		return false, fmt.Errorf("read the magic: %w", err)
	}
	if n != 5 {
		return false, fmt.Errorf("too few bytes read: %d", n)
	}

	if !bytes.Equal(magic[:], []byte("%PDF-")) {
		return false, nil
	}
	return true, nil
}

type InfoResponse struct {
	Version   string
	Languages []string
}

func Info() (InfoResponse, error) {
	languages, err := ocr.GetAvailableLanguages()
	if err != nil {
		return InfoResponse{}, err
	}
	return InfoResponse{
		Version:   ocr.Version(),
		Languages: languages,
	}, nil
}
