package main

import (
	"os"

	"github.com/Peltoche/zapette/cmd/zapette/commands"
	"github.com/Peltoche/zapette/internal/tools/buildinfos"
	"github.com/spf13/cobra"
)

const binaryName = "zapette"

type exitCode int

const (
	exitOK    exitCode = 0
	exitError exitCode = 1
)

func main() {
	code := mainRun()
	os.Exit(int(code))
}

func mainRun() exitCode {
	cmd := &cobra.Command{
		Use:     binaryName,
		Short:   "Observe and manage your server.",
		Version: buildinfos.Version(),
	}

	// Subcommands
	cmd.AddCommand(commands.NewRunCmd(binaryName))

	err := cmd.Execute()
	if err != nil {
		return exitError
	}

	return exitOK
}
