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

func ecloudFloatingIPRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "floatingip",
		Short: "sub-commands relating to floating IPs",
	}

	// Child commands
	cmd.AddCommand(ecloudFloatingIPListCmd(f))
	cmd.AddCommand(ecloudFloatingIPShowCmd(f))

	return cmd
}

func ecloudFloatingIPListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists floating IPs",
		Long:    "This command lists floating IPs",
		Example: "ukfast ecloud floatingip list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudFloatingIPList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudFloatingIPList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	fips, err := service.GetFloatingIPs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving floating IPs: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFloatingIPsProvider(fips))
}

func ecloudFloatingIPShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <floatingip: id>...",
		Short:   "Shows a floating IP",
		Long:    "This command shows one or more floating IPs",
		Example: "ukfast ecloud floatingip show fip-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing floating IP")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudFloatingIPShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudFloatingIPShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var fips []ecloud.FloatingIP
	for _, arg := range args {
		fip, err := service.GetFloatingIP(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving floating IP [%s]: %s", arg, err)
			continue
		}

		fips = append(fips, fip)
	}

	return output.CommandOutput(cmd, OutputECloudFloatingIPsProvider(fips))
}
