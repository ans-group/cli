package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudDatastoreRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datastore",
		Short: "sub-commands relating to datastores",
	}

	// Child commands
	cmd.AddCommand(ecloudDatastoreListCmd())
	cmd.AddCommand(ecloudDatastoreShowCmd())

	return cmd
}

func ecloudDatastoreListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists datastores",
		Long:    "This command lists datastores",
		Example: "ukfast ecloud datastore list",
		Run: func(cmd *cobra.Command, args []string) {
			ecloudDatastoreList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudDatastoreList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	datastores, err := service.GetDatastores(params)
	if err != nil {
		output.Fatalf("Error retrieving datastores: %s", err)
		return
	}

	outputECloudDatastores(datastores)
}

func ecloudDatastoreShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <datastore: id>...",
		Short:   "Shows a datastore",
		Long:    "This command shows one or more datastores",
		Example: "ukfast ecloud vm datastore 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing datastore")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudDatastoreShow(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudDatastoreShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	var datastores []ecloud.Datastore
	for _, arg := range args {
		datastoreID, err := strconv.Atoi(arg)
		if err != nil {
			OutputWithErrorLevelf("Invalid datastore ID [%s]", arg)
			continue
		}

		datastore, err := service.GetDatastore(datastoreID)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving datastore [%s]: %s", arg, err)
			continue
		}

		datastores = append(datastores, datastore)
	}

	outputECloudDatastores(datastores)
}
