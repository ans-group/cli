package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

func ecloudSiteRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "sub-commands relating to sites",
	}

	// Child commands
	cmd.AddCommand(ecloudSiteListCmd(f))
	cmd.AddCommand(ecloudSiteShowCmd(f))

	return cmd
}

func ecloudSiteListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists sites",
		Long:    "This command lists sites",
		Example: "ans ecloud site list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudSiteList(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("state", "", "state for filtering")

	return cmd
}

func ecloudSiteList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("state", "state"))
	if err != nil {
		return err
	}

	sites, err := service.GetSites(params)
	if err != nil {
		return fmt.Errorf("error retrieving sites: %s", err)
	}

	return output.CommandOutput(cmd, SiteCollection(sites))
}

func ecloudSiteShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <site: id>...",
		Short:   "Shows a site",
		Long:    "This command shows one or more sites",
		Example: "ans ecloud vm site 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing site")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudSiteShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudSiteShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var sites []ecloud.Site
	for _, arg := range args {
		siteID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid site ID [%s]", arg)
			continue
		}

		site, err := service.GetSite(siteID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving site [%s]: %s", arg, err)
			continue
		}

		sites = append(sites, site)
	}

	return output.CommandOutput(cmd, SiteCollection(sites))
}
