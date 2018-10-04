package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace"
	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/plugin/models"
)

const pluginName = "birds-eye"

type BirdsEye struct {
	UI terminal.UI
}

type Org struct {
	Resources []OrgResources `json:"resources"`
}

type OrgResources struct {
	Entity   OrgEntity   `json:"entity"`
	Metadata OrgMetadata `json:"metadata"`
}

type OrgEntity struct {
	Name             string `json:"name"`
	OrganizationGUID int    `json:"organization_guid"`
	AppsURL          string `json:"apps_url"`
}

type OrgMetadata struct {
	GUID string `json:"guid"`
}

type Space struct {
	Resources []SpaceResources `json:"resources"`
}

type SpaceResources struct {
	Entity SpaceEntity `json:"entity"`
}

type SpaceEntity struct {
	Name string `json:"name"`
}

func (c *BirdsEye) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == pluginName {

		var (
			err        error
			isLoggedIn bool
			orgs       []plugin_models.GetOrgs_Model
			spaces     []plugin_models.GetSpaces_Model
			orgNames   []string
		)

		ui := terminal.NewUI(
			os.Stdin,
			os.Stdout,
			terminal.NewTeePrinter(os.Stdout),
			trace.NewLogger(os.Stdout, false, "false", ""),
		)
		c.UI = ui

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

		for _, org := range orgs {
			orgNames = append(orgNames, org.Name)
		}

		table := ui.Table([]string{"", "", ""})
		table.Add("Orgs", "Spaces", "Apps")

		for _, o := range orgNames {
			table.Add(string(o))
		}
		err = table.Print()

		for _, org := range orgs {
			var orgResult Org
			url := fmt.Sprintf("/v2/organizations/%s/spaces", org.Guid)
			orgResult = c.UnmarshalOrg(url, cliConnection)

			var orgSpaces []string
			for _, space := range orgResult.Resources {
				orgSpaces = append(orgSpaces, space.Entity.Name)
			}
			fmt.Print("\n\n")
			fmt.Printf("All spaces in %s:\n\n%s", org.Name, strings.Join(orgSpaces, "\n"))

			for _, space := range orgResult.Resources {
				var spaceResult Space

				getSpaceAppsRequest := fmt.Sprintf("/v2/spaces/%s/apps", space.Metadata.GUID)
				spaceResult = c.UnmarshalSpace(getSpaceAppsRequest, cliConnection)

				var appsInSpace []string
				displaySpaceApps(appsInSpace, spaceResult, space)
			}
		}
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

func displaySpaceApps(spaceApps []string, space Space, sResult OrgResources) {
	for _, app := range space.Resources {
		spaceApps = append(spaceApps, app.Entity.Name)
	}
	fmt.Print("\n\n")
	fmt.Printf("All apps in %s:\n\n%s", sResult.Entity.Name, strings.Join(spaceApps, "\n"))
}

func (c BirdsEye) UnmarshalOrg(url string, cliConnection plugin.CliConnection) Org {
	var orgResult Org
	cmd := []string{"curl", url}

	result, _ := cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(result, "")), &orgResult)

	return orgResult
}

func (c BirdsEye) UnmarshalSpace(url string, cliConnection plugin.CliConnection) Space {
	var spaceResult Space
	cmd := []string{"curl", url}

	result, _ := cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(result, "")), &spaceResult)

	return spaceResult
}

func main() {
	plugin.Start(new(BirdsEye))
}
