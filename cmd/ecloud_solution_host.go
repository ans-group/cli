package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSolutionHostRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host",
		Short: "sub-commands relating to solution hosts",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionHostListCmd())

	return cmd
}

func ecloudSolutionHostListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists solution hosts",
		Long:    "This command lists solution hosts",
		Example: "ukfast ecloud solution host list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionHostList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionHostList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
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

	hosts, err := service.GetSolutionHosts(solutionID, params)
	if err != nil {
		output.Fatalf("Error retrieving solution hosts: %s", err)
		return
	}

	outputECloudHosts(hosts)
}
