package ecloudv2

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ECloudV2RootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ecloud",
		Short: "Commands relating to eCloud service",
	}

	// Child root commands
	cmd.AddCommand(ecloudVPCRootCmd(f))
	cmd.AddCommand(ecloudInstanceRootCmd(f))
	cmd.AddCommand(ecloudFloatingIPRootCmd(f))
	cmd.AddCommand(ecloudFirewallRuleRootCmd(f))
	cmd.AddCommand(ecloudRegionRootCmd(f))

	return cmd
}

type ecloudServiceCobraRunEFunc func(service ecloud.ECloudService, cmd *cobra.Command, args []string) error

func ecloudCobraRunEFunc(f factory.ClientFactory, rf ecloudServiceCobraRunEFunc) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		c, err := f.NewClient()
		if err != nil {
			return err
		}

		return rf(c.ECloudService(), cmd, args)
	}
}
