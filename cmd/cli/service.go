package cli

import (
	gl "github.com/rafa-mori/goforge/logger"
	"github.com/spf13/cobra"
)

func ServiceCmdList() []*cobra.Command {
	return []*cobra.Command{
		startCommand(),
	}
}

func startCommand() *cobra.Command {
	var debug bool

	var startCmd = &cobra.Command{
		Use: "start",
		Annotations: GetDescriptions([]string{
			"Start some command.",
			"This command is used to start the GoForge service with the specified configuration.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			if debug {
				gl.SetDebug(true)
				gl.Log("debug", "Debug mode enabled")
			}
			gl.Log("success", "GoForge service started successfully")
		},
	}

	startCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")

	return startCmd
}
