//go:build mage

package main

import (
	"context"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Mod mg.Namespace

// mod:download runs go mod download.
func (Mod) Download(ctx context.Context) error {
	return sh.RunV("go", "mod", "download")
}

func Build(ctx context.Context) error {
	mg.Deps(Mod.Download)
	if err := sh.RunV("go", "build", "github.com/gomoni/kmd/cmd/kmc"); err != nil {
		return err
	}
	if err := sh.RunV("go", "build", "github.com/gomoni/kmd/cmd/kmd"); err != nil {
		return err
	}
	return nil
}

type Test mg.Namespace

// test:unit runs the unit tests.
func (Test) Unit(ctx context.Context) error {
	if err := sh.RunV("go", "test", "-v", "./..."); err != nil {
		return err
	}
	return nil
}
