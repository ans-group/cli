package ecloud_v2

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/ptr"
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
	cmd.AddCommand(ecloudRouterCreateCmd(f))
	cmd.AddCommand(ecloudRouterUpdateCmd(f))
	cmd.AddCommand(ecloudRouterDeleteCmd(f))

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

func ecloudRouterCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a router",
		Long:    "This command creates a router",
		Example: "ukfast ecloud router create --vpc vpc-abcdef12 --az az-abcdef12",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudRouterCreate(c.ECloudService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of router")
	cmd.Flags().String("vpc", "", "ID of VPC")
	cmd.MarkFlagRequired("vpc")
	cmd.Flags().String("az", "", "ID of availability zone")
	cmd.MarkFlagRequired("az")

	return cmd
}

func ecloudRouterCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateRouterRequest{}
	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		createRequest.Name = ptr.String(name)
	}
	createRequest.VPCID, _ = cmd.Flags().GetString("vpc")
	createRequest.AvailabilityZoneID, _ = cmd.Flags().GetString("az")

	routerID, err := service.CreateRouter(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating router: %s", err)
	}

	router, err := service.GetRouter(routerID)
	if err != nil {
		return fmt.Errorf("Error retrieving new router: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudRoutersProvider([]ecloud.Router{router}))
}

func ecloudRouterUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <router: id>...",
		Short:   "Updates a router",
		Long:    "This command updates one or more routers",
		Example: "ukfast ecloud router update rtr-abcdef12 --name \"my router\"",
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

			return ecloudRouterUpdate(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name of router")

	return cmd
}

func ecloudRouterUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchRouterRequest{}

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		patchRequest.Name = ptr.String(name)
	}

	var routers []ecloud.Router
	for _, arg := range args {
		err := service.PatchRouter(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating router [%s]: %s", arg, err)
			continue
		}

		router, err := service.GetRouter(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated router [%s]: %s", arg, err)
			continue
		}

		routers = append(routers, router)
	}

	return output.CommandOutput(cmd, OutputECloudRoutersProvider(routers))
}

func ecloudRouterDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <router: id...>",
		Short:   "Removes a router",
		Long:    "This command removes one or more routers",
		Example: "ukfast ecloud router delete rtr-abcdef12",
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

			ecloudRouterDelete(c.ECloudService(), cmd, args)
			return nil
		},
	}
}

func ecloudRouterDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		err := service.DeleteRouter(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing router [%s]: %s", arg, err)
		}
	}
}
