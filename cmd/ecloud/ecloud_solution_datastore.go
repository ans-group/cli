package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	flaghelper "github.com/ukfast/cli/internal/pkg/helper/flag"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSolutionDatastoreRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datastore",
		Short: "sub-commands relating to solution datastores",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionDatastoreListCmd(f))

	return cmd
}

func ecloudSolutionDatastoreListCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudSolutionDatastoreList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionDatastoreList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid solution ID [%s]", args[0])
	}

	params, err := flaghelper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	datastores, err := service.GetSolutionDatastores(solutionID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving solution datastores: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudDatastoresProvider(datastores))
}
