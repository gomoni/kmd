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
	debug bool

	client client.HTTP
	ctx    context.Context
}

func (g Globals) Debug(msg string, args ...any) {
	if g.debug {
		slog.Debug(msg, args...)
	}
}

type InfoCmd struct{}

func (r *InfoCmd) Run(g Globals) error {
	g.Debug("Running Info")
	return g.client.Info(g.ctx, os.Stdout)
}

type CLI struct {
	Debug bool `help:"Enable debug mode."`

	Info InfoCmd `cmd:"" help:"show the server information"`
}

func maine() error {
	var cli CLI
	parsed := kong.Parse(&cli)
	if cli.Debug {
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
		debug:  cli.Debug,
		client: client,
		ctx:    context.Background(),
	})
}

func main() {
	if err := maine(); err != nil {
		log.Fatal(err)
	}
}
