package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
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
		Example: "ans ecloud instance nic list i-abcdef12",
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
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	nics, err := service.GetInstanceNICs(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving instance NICs: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudNICsProvider(nics))
}
