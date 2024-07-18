package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
)

type HTTP struct {
	prefix string
	client *http.Client
}

func NewHTTP(prefix string) HTTP {
	return HTTP{
		client: http.DefaultClient,
		prefix: prefix,
	}
}

func NewUnix(path string) (HTTP, error) {
	unixClient := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", path)
			},
		},
	}
	return NewHTTP("http://unix").WithHTTPClient(&unixClient), nil
}

func (c HTTP) WithHTTPClient(client *http.Client) HTTP {
	ret := c
	ret.client = client
	return ret
}

func (c HTTP) OCR(ctx context.Context, w io.Writer, r io.Reader) error {
	if r == nil {
		return fmt.Errorf("input is empty")
	}

	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	fw, err := mw.CreateFormFile("file", "a.file")
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}
	var buf [4096]byte
	_, err = io.CopyBuffer(fw, r, buf[:])
	if err != nil {
		return fmt.Errorf("read from the file: %w", err)
	}
	// TODO: implement the support for OCR params
	/*
		var errs []error
			addError := func(err error) {
				if err != nil {
					errs = append(errs, err)
				}
			}
				if params.Languages != nil {
					err = mw.WriteField("languages", strings.Join(params.Languages, ","))
					addError(err)
				}
				if params.Whitelist != "" {
					err = mw.WriteField("whitelist", params.Whitelist)
					addError(err)
				}
				if params.HOCR {
					err = mw.WriteField("format", "hocr")
					addError(err)
				}
		if errs != nil {
			return fmt.Errorf("mutlipart errors: %w", errors.Join(errs...))
		}
	*/
	mw.Close()

	req, err := postPlain(ctx, c.prefix+"/ocr", &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return c.do(w, req)
}

func (c HTTP) Info(ctx context.Context, w io.Writer) error {
	req, err := getPlain(ctx, c.prefix+"/")
	if err != nil {
		return err
	}
	return c.do(w, req)
}

func (c HTTP) do(w io.Writer, r *http.Request) error {
	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("http response %d: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("http response %d: %s", resp.StatusCode, string(b))
	}
	var buf [1024]byte
	_, err = io.CopyBuffer(w, resp.Body, buf[:])
	if err != nil {
		return err
	}
	return nil
}

func postPlain(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		body,
	)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("content-type", "text/plain")
	return req, nil
}

func getPlain(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("content-type", "text/plain")
	return req, nil
}
