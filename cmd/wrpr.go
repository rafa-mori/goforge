package main

import (
	cc "github.com/rafa-mori/goforge/cmd/cli"
	gl "github.com/rafa-mori/goforge/logger"
	vs "github.com/rafa-mori/goforge/version"
	"github.com/spf13/cobra"

	"os"
	"strings"
)

type GoForge struct {
	parentCmdName string
	printBanner   bool
}

func (m *GoForge) Alias() string {
	return ""
}
func (m *GoForge) ShortDescription() string {
	return "GoForge is a minimalistic backend service with Go."
}
func (m *GoForge) LongDescription() string {
	return `GoForge: A minimalistic backend service with Go.`
}
func (m *GoForge) Usage() string {
	return "article [command] [args]"
}
func (m *GoForge) Examples() []string {
	return []string{"article some-command",
		"article another-command --option value",
		"article yet-another-command --flag"}
}
func (m *GoForge) Active() bool {
	return true
}
func (m *GoForge) Module() string {
	return "article"
}
func (m *GoForge) Execute() error {
	return m.Command().Execute()
}
func (m *GoForge) Command() *cobra.Command {
	gl.Log("debug", "Starting GoForge CLI...")

	var rtCmd = &cobra.Command{
		Use:     m.Module(),
		Aliases: []string{m.Alias()},
		Example: m.concatenateExamples(),
		Version: vs.GetVersion(),
		Annotations: cc.GetDescriptions([]string{
			m.LongDescription(),
			m.ShortDescription(),
		}, m.printBanner),
	}

	rtCmd.AddCommand(cc.ServiceCmdList()...)
	rtCmd.AddCommand(vs.CliCommand())

	// Set usage definitions for the command and its subcommands
	setUsageDefinition(rtCmd)
	for _, c := range rtCmd.Commands() {
		setUsageDefinition(c)
		if !strings.Contains(strings.Join(os.Args, " "), c.Use) {
			if c.Short == "" {
				c.Short = c.Annotations["description"]
			}
		}
	}

	return rtCmd
}
func (m *GoForge) SetParentCmdName(rtCmd string) {
	m.parentCmdName = rtCmd
}
func (m *GoForge) concatenateExamples() string {
	examples := ""
	rtCmd := m.parentCmdName
	if rtCmd != "" {
		rtCmd = rtCmd + " "
	}
	for _, example := range m.Examples() {
		examples += rtCmd + example + "\n  "
	}
	return examples
}
func RegX() *GoForge {
	var printBannerV = os.Getenv("GOFORGE_PRINT_BANNER")
	if printBannerV == "" {
		printBannerV = "true"
	}

	return &GoForge{
		printBanner: strings.ToLower(printBannerV) == "true",
	}
}
