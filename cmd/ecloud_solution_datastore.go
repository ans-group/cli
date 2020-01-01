package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSolutionDatastoreRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datastore",
		Short: "sub-commands relating to solution datastores",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionDatastoreListCmd())

	return cmd
}

func ecloudSolutionDatastoreListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists solution datastores",
		Long:    "This command lists solution datastores",
		Example: "ukfast ecloud solution datastore list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionDatastoreList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionDatastoreList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
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

	datastores, err := service.GetSolutionDatastores(solutionID, params)
	if err != nil {
		output.Fatalf("Error retrieving solution datastores: %s", err)
		return
	}

	outputECloudDatastores(datastores)
}
