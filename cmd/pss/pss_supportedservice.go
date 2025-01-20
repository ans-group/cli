package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssSupportedServiceRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supportedservice",
		Short: "sub-commands relating to supported services",
	}

	// Child commands
	cmd.AddCommand(pssSupportedServiceListCmd(f))

	return cmd
}

func pssSupportedServiceListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists supported services",
		Long:    "This command lists supported services",
		Example: "ans pss supportedservice list",
		RunE:    pssCobraRunEFunc(f, pssSupportedServiceList),
	}
}

func pssSupportedServiceList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	problems, err := service.GetSupportedServices(params)
	if err != nil {
		return err
	}

	return output.CommandOutput(cmd, SupportedServiceCollection(problems))
}
