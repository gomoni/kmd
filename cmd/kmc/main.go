package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gomoni/kmd/internal/client"
	kos "github.com/gomoni/kmd/internal/os"

	"github.com/alecthomas/kong"
)

type Globals struct {
	// params
	verbose bool

	// other global stuff
	client client.HTTP
	ctx    context.Context
}

func (g Globals) Debug(msg string, args ...any) {
	if g.verbose {
		slog.Debug(msg, args...)
	}
}

type InfoCmd struct{}

func (r *InfoCmd) Run(g Globals) error {
	g.Debug("Running Info")
	return g.client.Info(g.ctx, os.Stdout)
}

type OCRCmd struct {
	Languages []string `optional:"" short:"l" help:"languages to use for OCR (default to eng(lish))"`
	File      string   `arg:"" required:"" help:"file to process"`
}

func (r *OCRCmd) Run(g Globals) error {
	g.Debug("Running OCR", "file", r.File, "languages", r.Languages)
	f, err := os.Open(r.File)
	if err != nil {
		return fmt.Errorf("os.Open(%s): %w", r.File, err)
	}
	defer f.Close()

	params := client.OCRParams{
		Languages: r.Languages,
	}
	return g.client.OCR(g.ctx, os.Stdout, f, params)
}

type CLI struct {
	Verbose bool   `help:"Enable verbose mode."`
	Host    string `optional:"" short:"H" env:"KMD_HOST" help:"Server host, defaults to $XDG_RUNTIME_DIR/kmd.sock or to KMD_HOST environment variable."`

	Info InfoCmd `cmd:"" default:"1" help:"show the server information"`
	OCR  OCRCmd  `cmd:"" help:"extract the text from the input image or pdf file"`
}

func maine() error {
	var cli CLI
	parsed := kong.Parse(&cli)
	if cli.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	if cli.Host == "" {
		cli.Host = kos.DefaultSocketPath()
	}

	client, err := client.NewUnix(string(cli.Host))
	if err != nil {
		return fmt.Errorf("client.NewUnix(%s): %w", cli.Host, err)
	}

	return parsed.Run(&cli, Globals{
		verbose: cli.Verbose,
		client:  client,
		ctx:     context.Background(),
	})
}

func main() {
	if err := maine(); err != nil {
		log.Fatal(err)
	}
}
