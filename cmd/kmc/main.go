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
	verbose bool

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
	File string `arg:"" required:"" help:"file to process"`
}

func (r *OCRCmd) Run(g Globals) error {
	g.Debug("Running OCR", "file", r.File)
	f, err := os.Open(r.File)
	if err != nil {
		return fmt.Errorf("os.Open(%s): %w", r.File, err)
	}
	defer f.Close()

	return g.client.OCR(g.ctx, os.Stdout, f)
}

type CLI struct {
	Verbose bool `help:"Enable verbose mode."`

	Info InfoCmd `cmd:"" help:"show the server information"`
	OCR  OCRCmd  `cmd:"" help:"extract the text from the input image or pdf file"`
}

func maine() error {
	var cli CLI
	parsed := kong.Parse(&cli)
	if cli.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	// TODO: KMD_HOST support with kong
	// TODO: --host/-H docker-like support
	path := kos.DefaultSocketPath()

	client, err := client.NewUnix(path)
	if err != nil {
		return fmt.Errorf("client.NewUnix(%s): %w", path, err)
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
