package ecloud_v2

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudRouterRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "router",
		Short: "sub-commands relating to routers",
	}

	// Child commands
	cmd.AddCommand(ecloudRouterListCmd(f))
	cmd.AddCommand(ecloudRouterShowCmd(f))

	return cmd
}

func ecloudRouterListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists routers",
		Long:    "This command lists routers",
		Example: "ukfast ecloud router list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudRouterList(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Router name for filtering")

	return cmd
}

func ecloudRouterList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	helper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, helper.NewStringFilterFlag("name", "name"))

	routers, err := service.GetRouters(params)
	if err != nil {
		return fmt.Errorf("Error retrieving routers: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudRoutersProvider(routers))
}

func ecloudRouterShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <router: id>...",
		Short:   "Shows a router",
		Long:    "This command shows one or more routers",
		Example: "ukfast ecloud router show rtr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing router")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudRouterShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudRouterShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var routers []ecloud.Router
	for _, arg := range args {
		router, err := service.GetRouter(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving router [%s]: %s", arg, err)
			continue
		}

		routers = append(routers, router)
	}

	return output.CommandOutput(cmd, OutputECloudRoutersProvider(routers))
}
