package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"
var buildDate = "notset"
var gitHash = ""

func init() {
	rootCmd.AddCommand(versionCmd)

	rootCmd.Version = fmt.Sprintf("%s [%s] (%s)", version, gitHash, buildDate)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run:   versionCommand,
}

func versionCommand(cmd *cobra.Command, args []string) {
	fmt.Printf("%s version %s [%s] (%s)\n", rootCmd.Name(), version, gitHash, buildDate)
}
