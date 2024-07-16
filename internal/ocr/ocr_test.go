package ocr_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	. "github.com/gomoni/kmd/internal/ocr"
	"github.com/gomoni/kmd/internal/render"

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

func TestClientImageReader(t *testing.T) {
	t.Parallel()

	f, err := os.Open("../testdata/hello.png")
	require.NoError(t, err)
	cleanupe(t, f.Close)
	c := NewClient()
	t.Cleanup(c.Close)
	err = c.ImageReader(f)
	require.NoError(t, err)
	text, err := c.Text()
	require.NoError(t, err)
	require.Equal(t, "Hello, world!", text)
}

func TestSmartClient(t *testing.T) {
	t.Parallel()
	var testCases = []string{"../testdata/hello.png", "../testdata/hello.pdf"}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			t.Parallel()

			f, err := os.Open(tc)
			require.NoError(t, err)
			cleanupe(t, f.Close)
			st, err := f.Stat()
			require.NoError(t, err)
			t.Logf("%s size: %d", tc, st.Size())
			client := NewClient()
			sc := NewSmartClient(client, pool)
			t.Cleanup(client.Close)
			err = sc.ImageReader(f)
			require.NoError(t, err)
			text, err := sc.Text()
			require.NoError(t, err)
			require.Equal(t, "Hello, world!", text)
		})
	}
}

func TestIsPDF(t *testing.T) {
	f, err := os.Open("../testdata/hello.pdf")
	require.NoError(t, err)
	cleanupe(t, f.Close)

	is, err := IsPDF(f)
	require.NoError(t, err)
	require.True(t, is)
}

func TestIsNotPDF(t *testing.T) {
	f, err := os.Open("../testdata/hello.png")
	require.NoError(t, err)
	cleanupe(t, f.Close)

	is, err := IsPDF(f)
	require.NoError(t, err)
	require.False(t, is)
}

func cleanupe(t *testing.T, fun func() error) {
	t.Helper()
	t.Cleanup(func() { require.NoError(t, fun()) })
}
