package server_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gomoni/kmd/internal/client"
	"github.com/gomoni/kmd/internal/render"
	. "github.com/gomoni/kmd/internal/server"

	"github.com/stretchr/testify/require"
)

var pool render.Pool

func TestMain(m *testing.M) {
	var err error
	t0 := time.Now()
	pool, err = render.NewPool()
	if err != nil {
		log.Fatalf("NewPool err: %s", err)
	}
	t1 := time.Now()

	fmt.Printf("pdfium pool created in %+v\n", t1.Sub(t0))

	ret := m.Run()

	err = pool.Close()
	if err != nil {
		log.Fatalf("pool.Close err: %s", err)
	}

	os.Exit(ret)
}

func TestHTTPHandler(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	handler := NewOCR(32<<20, pool)
	server := httptest.NewServer(handler)

	var testCases = []string{"../testdata/hello.png", "../testdata/hello.pdf"}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			t.Parallel()

			r, err := os.Open(tc)
			require.NoError(t, err)
			require.NotNil(t, r)
			t.Cleanup(
				func() {
					require.NoError(t, r.Close())
				})

			params := client.OCRParams{Languages: []string{"eng"}}
			client := client.NewHTTP(server.URL).WithHTTPClient(server.Client())
			var out bytes.Buffer
			err = client.OCR(ctx, &out, r, params)
			require.NoError(t, err)

			require.Equal(t, "Hello, world!", out.String())
		})
	}
}
