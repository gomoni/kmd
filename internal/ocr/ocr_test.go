package ocr_test

import (
	"embed"
	"testing"

	. "github.com/gomoni/kmd/internal/ocr"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/hello.png
var helloPng []byte

//go:embed testdata/hello.png
var helloPngFile embed.FS

func TestClientImagePath(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		scenario string
		given    func(c Client) error
	}{
		{
			scenario: "ImagePath",
			given: func(c Client) error {
				return c.ImagePath("testdata/hello.png")
			},
		},
		{
			scenario: "ImageBytes",
			given: func(c Client) error {
				return c.ImageBytes(helloPng)
			},
		},
		{
			scenario: "ImageReader",
			given: func(c Client) error {
				r, err := helloPngFile.Open("testdata/hello.png")
				if err != nil {
					return err
				}
				return c.ImageReader(r)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			t.Parallel()
			c := NewClient()
			t.Cleanup(c.Close)
			err := tt.given(c)
			require.NoError(t, err)
			text, err := c.Text()
			require.NoError(t, err)
			require.Equal(t, "Hello, world!", text)
		})
	}
}

func TestClientImageBytes(t *testing.T) {
	c := NewClient()
	t.Cleanup(c.Close)
	err := c.ImageBytes(helloPng)
	require.NoError(t, err)
	text, err := c.Text()
	require.NoError(t, err)
	require.Equal(t, "Hello, world!", text)
}
