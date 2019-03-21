package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSiteRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "sub-commands relating to sites",
	}

	// Child commands
	cmd.AddCommand(ecloudSiteListCmd())
	cmd.AddCommand(ecloudSiteShowCmd())

	return cmd
}

func ecloudSiteListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists sites",
		Long:    "This command lists sites",
		Example: "ukfast ecloud site list",
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSiteList(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("state", "", "state for filtering")

	return cmd
}

func ecloudSiteList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	if cmd.Flags().Changed("state") {
		filterState, _ := cmd.Flags().GetString("state")
		params.WithFilter(helper.GetFilteringInferOperator("state", filterState))
	}

	sites, err := service.GetSites(params)
	if err != nil {
		output.Fatalf("Error retrieving sites: %s", err)
		return
	}

	outputECloudSites(sites)
}

func ecloudSiteShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <site: id>...",
		Short:   "Shows a site",
		Long:    "This command shows one or more sites",
		Example: "ukfast ecloud vm site 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing site")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSiteShow(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSiteShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	var sites []ecloud.Site
	for _, arg := range args {
		siteID, err := strconv.Atoi(arg)
		if err != nil {
			OutputWithErrorLevelf("Invalid site ID [%s]", arg)
			continue
		}

		site, err := service.GetSite(siteID)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving site [%s]: %s", arg, err)
			continue
		}

		sites = append(sites, site)
	}

	outputECloudSites(sites)
}
