package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSolutionSiteRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "sub-commands relating to solution sites",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionSiteListCmd(f))

	return cmd
}

func ecloudSolutionSiteListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists solution sites",
		Long:    "This command lists solution sites",
		Example: "ukfast ecloud solution site list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionSiteList(f.NewClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionSiteList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid solution ID [%s]", args[0])
	}

	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	sites, err := service.GetSolutionSites(solutionID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving solution sites: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudSitesProvider(sites))
}
