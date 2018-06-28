package main

import (
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/plugin/models"
)

const pluginName = "birds-eye"

type BirdsEye struct{}

func (c *BirdsEye) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == pluginName {

		var (
			err error

			isLoggedIn bool
			orgs       []plugin_models.GetOrgs_Model
			spaces     []plugin_models.GetSpaces_Model
			apps       []plugin_models.GetAppsModel

			orgNames []string
			// spaceNames []string
			appNames []string

			orgGUIDs []string
		)

		if _, err = cliConnection.HasAPIEndpoint(); err != nil {
			fmt.Println("No API endpoint set")
		}

		if isLoggedIn, err = cliConnection.IsLoggedIn(); err != nil {
			fmt.Printf("Logged in? %t", isLoggedIn)
		}

		if orgs, err = cliConnection.GetOrgs(); err != nil {
			fmt.Printf("Error getting orgs: %v", orgs)
		}

		if spaces, err = cliConnection.GetSpaces(); err != nil {
			fmt.Printf("Error getting spaces: %v", spaces)
		}

		if apps, err = cliConnection.GetApps(); err != nil {
			fmt.Printf("Error getting apps: %v", apps)
		}

		for _, org := range orgs {
			orgNames = append(orgNames, org.Name)
			orgGUIDs = append(orgGUIDs, org.Guid)
		}

		fmt.Print("All orgs:\n\n", strings.Join(orgNames, "\n"))

		url := fmt.Sprintf("/v2/organizations/%s/spaces", orgGUIDs[0])

		orgSpaces, err := cliConnection.CliCommandWithoutTerminalOutput("curl", url)
		if err != nil {
			fmt.Printf("Error getting spaces from org: %s", orgNames[0])
		}

		fmt.Print("All spaces:\n\n", strings.Join(orgSpaces, "\n"))
		fmt.Print("All apps:\n\n", strings.Join(appNames, "\n"))
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
