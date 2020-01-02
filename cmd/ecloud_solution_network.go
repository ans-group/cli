package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSolutionNetworkRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "sub-commands relating to solution networks",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionNetworkListCmd())

	return cmd
}

func ecloudSolutionNetworkListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists solution networks",
		Long:    "This command lists solution networks",
		Example: "ukfast ecloud solution network list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionNetworkList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionNetworkList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
	}

	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	networks, err := service.GetSolutionNetworks(solutionID, params)
	if err != nil {
		output.Fatalf("Error retrieving solution networks: %s", err)
		return
	}

	outputECloudNetworks(networks)
}
