package ecloud

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudInstanceFloatingIPRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "floatingip",
		Short: "sub-commands relating to instance floating IPs",
	}

	// Child commands
	cmd.AddCommand(ecloudInstanceFloatingIPListCmd(f))

	return cmd
}

func ecloudInstanceFloatingIPListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists instance floating IPs",
		Long:    "This command lists instance floating IPs",
		Example: "ukfast ecloud instance floatingip list i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceFloatingIPList),
	}

	cmd.Flags().String("name", "", "Floating IP name for filtering")

	return cmd
}

func ecloudInstanceFloatingIPList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	fips, err := service.GetInstanceFloatingIPs(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving instance floating IPs: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFloatingIPsProvider(fips))
}
