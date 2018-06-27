package main

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
)

const pluginName = "birds-eye"

type BirdsEye struct{}

func (c *BirdsEye) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == pluginName {
		fmt.Println("Running the command")
	}
}

func (c *BirdsEye) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: pluginName,
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     pluginName,
				HelpText: "Displays all orgs, spaces, and apps for the CF instance",
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s", pluginName),
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(BirdsEye))
}
