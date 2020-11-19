package ecloudv2

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudInstanceNICRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nic",
		Short: "sub-commands relating to instance NICs",
	}

	// Child commands
	cmd.AddCommand(ecloudInstanceNICListCmd(f))

	return cmd
}

func ecloudInstanceNICListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists instance nics",
		Long:    "This command lists instance nics",
		Example: "ukfast ecloud instance nic list i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceNICList),
	}

	cmd.Flags().String("name", "", "NIC name for filtering")

	return cmd
}

func ecloudInstanceNICList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	helper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, helper.NewStringFilterFlag("name", "name"))

	nics, err := service.GetInstanceNICs(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving instance NICs: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudNICsProvider(nics))
}
