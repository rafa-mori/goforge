package main

import (
	cc "github.com/faelmori/article/cmd/cli"
	gl "github.com/faelmori/article/logger"
	vs "github.com/faelmori/article/version"
	"github.com/spf13/cobra"

	"os"
	"strings"
)

type ArticleMod struct {
	parentCmdName string
	printBanner   bool
}

func (m *ArticleMod) Alias() string {
	return ""
}
func (m *ArticleMod) ShortDescription() string {
	return "ArticleMod is a minimalistic backend service with Go."
}
func (m *ArticleMod) LongDescription() string {
	return `ArticleMod: A minimalistic backend service with Go.`
}
func (m *ArticleMod) Usage() string {
	return "article [command] [args]"
}
func (m *ArticleMod) Examples() []string {
	return []string{"article some-command",
		"article another-command --option value",
		"article yet-another-command --flag"}
}
func (m *ArticleMod) Active() bool {
	return true
}
func (m *ArticleMod) Module() string {
	return "article"
}
func (m *ArticleMod) Execute() error {
	return m.Command().Execute()
}
func (m *ArticleMod) Command() *cobra.Command {
	gl.Log("debug", "Starting ArticleMod CLI...")

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
func (m *ArticleMod) SetParentCmdName(rtCmd string) {
	m.parentCmdName = rtCmd
}
func (m *ArticleMod) concatenateExamples() string {
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
func RegX() *ArticleMod {
	var printBannerV = os.Getenv("ARTICLE_PRINT_BANNER")
	if printBannerV == "" {
		printBannerV = "true"
	}

	return &ArticleMod{
		printBanner: strings.ToLower(printBannerV) == "true",
	}
}
