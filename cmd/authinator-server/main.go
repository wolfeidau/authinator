package main

import "github.com/spf13/cobra"

var (
	// Version The version of the application (set by make file)
	Version = "UNKNOWN"

	cmdRoot = &cobra.Command{
		Use:   "authinator-server",
		Short: "authinator server",
		Long:  ``,
	}

	rootOpts struct {
		Debug bool
	}
)

func init() {
	cmdRoot.PersistentFlags().BoolVar(&rootOpts.Debug, "debug", false, "Log debug information")
}

func main() {
	cmdRoot.Execute()
}
