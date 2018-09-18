package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/plugin/models"
)

const pluginName = "birds-eye"

type BirdsEye struct{}

type Org struct {
	Resources []OrgResources `json:"resources"`
}

type OrgResources struct {
	Entity OrgEntity `json:"entity"`
}

type OrgEntity struct {
	Name             string `json:"name"`
	OrganizationGUID int    `json:"organization_guid"`
	AppsURL          string `json:"apps_url"`
}

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
		}
		fmt.Print("All orgs:\n\n*s", strings.Join(orgNames, "\n"))

		for _, org := range orgs {
			var orgResult Org
			url := fmt.Sprintf("/v2/organizations/%s/spaces", org.Guid)
			orgResult = c.UnmarshalOrg(url, cliConnection)

			var orgSpaces, appsURLsInSpace []string
			for _, space := range orgResult.Resources {
				orgSpaces = append(orgSpaces, space.Entity.Name)
				appsURLsInSpace = append(appsURLsInSpace, space.Entity.AppsURL)
			}

			fmt.Print("\n\n")
			fmt.Printf("All spaces in %s:\n\n* %s", org.Name, strings.Join(orgSpaces, "\n"))

			for _, app := range apps {
				appNames = append(appNames, app.Name)
			}

			// fmt.Print("\n\n")
			// fmt.Print("All apps:\n\n", strings.Join(appNames, "\n"))
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

func (c BirdsEye) UnmarshalOrg(url string, cliConnection plugin.CliConnection) Org {
	var orgResult Org
	cmd := []string{"curl", url}

	result, _ := cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(result, "")), &orgResult)

	return orgResult
}

func main() {
	plugin.Start(new(BirdsEye))
}
