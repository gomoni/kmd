package server_test

import (
	"bytes"
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gomoni/kmd/internal/client"
	. "github.com/gomoni/kmd/internal/server"

	"github.com/stretchr/testify/require"
)

func TestHTTPHandler(t *testing.T) {
	ctx := context.Background()

	r, err := os.Open("../ocr/testdata/hello.png")
	require.NoError(t, err)
	require.NotNil(t, r)
	t.Cleanup(
		func() {
			require.NoError(t, r.Close())
		})

	handler := NewOCR(32 << 20)
	server := httptest.NewServer(handler)

	client := client.NewHTTP(server.URL).WithHTTPClient(server.Client())
	var out bytes.Buffer
	err = client.OCR(ctx, &out, r)
	require.NoError(t, err)

	require.Equal(t, "Hello, world!", out.String())
}
