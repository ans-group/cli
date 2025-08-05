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

func ecloudDatastoreRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datastore",
		Short: "sub-commands relating to datastores",
	}

	// Child commands
	cmd.AddCommand(ecloudDatastoreListCmd(f))
	cmd.AddCommand(ecloudDatastoreShowCmd(f))

	return cmd
}

func ecloudDatastoreListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists datastores",
		Long:    "This command lists datastores",
		Example: "ans ecloud datastore list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudDatastoreList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudDatastoreList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	datastores, err := service.GetDatastores(params)
	if err != nil {
		return fmt.Errorf("error retrieving datastores: %s", err)
	}

	return output.CommandOutput(cmd, DatastoreCollection(datastores))
}

func ecloudDatastoreShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <datastore: id>...",
		Short:   "Shows a datastore",
		Long:    "This command shows one or more datastores",
		Example: "ans ecloud vm datastore 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing datastore")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudDatastoreShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudDatastoreShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var datastores []ecloud.Datastore
	for _, arg := range args {
		datastoreID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid datastore ID [%s]", arg)
			continue
		}

		datastore, err := service.GetDatastore(datastoreID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving datastore [%s]: %s", arg, err)
			continue
		}

		datastores = append(datastores, datastore)
	}

	return output.CommandOutput(cmd, DatastoreCollection(datastores))
}
