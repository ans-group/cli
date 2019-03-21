package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSolutionSiteRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "sub-commands relating to solution sites",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionSiteListCmd())

	return cmd
}

func ecloudSolutionSiteListCmd() *cobra.Command {
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
			ecloudSolutionSiteList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionSiteList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
	}

	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	sites, err := service.GetSolutionSites(solutionID, params)
	if err != nil {
		output.Fatalf("Error retrieving solution sites: %s", err)
		return
	}

	outputECloudSites(sites)
}
