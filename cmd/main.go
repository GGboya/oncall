package main

import (
	"fmt"
	"os"

	"oncall/cmd/app"

	"github.com/sirupsen/logrus"
)

func main() {
	deps, err := app.NewAppDependencies()
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize app dependencies")
		os.Exit(1)
	}
	cmd := app.NewCommand(deps)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1) // nolint:gocritic
	}
}
