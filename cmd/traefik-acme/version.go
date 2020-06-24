package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nolint: gochecknoglobals // these have to be variables for the linker to change the values
var (
	version   = "dev"
	buildDate = "notset"
	gitHash   = ""
)

// nolint: gochecknoglobals // cobra uses globals in main
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version",
	Run:     versionCommand,
}

// nolint:gochecknoinits // init is used in main for cobra
func init() {
	rootCmd.AddCommand(versionCmd)

	rootCmd.Version = fmt.Sprintf("%s [%s] (%s)", version, gitHash, buildDate)
}

func versionCommand(cmd *cobra.Command, args []string) {
	fmt.Printf("%s version %s [%s] (%s)\n", rootCmd.Name(), version, gitHash, buildDate)
}
