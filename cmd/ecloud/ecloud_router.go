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
		RunE:    ecloudCobraRunEFunc(f, ecloudRouterList),
	}

	cmd.Flags().String("name", "", "Router name for filtering")

	return cmd
}

func ecloudRouterList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

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
		RunE: ecloudCobraRunEFunc(f, ecloudRouterShow),
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
		Short:   "Creates an router",
		Long:    "This command creates an router",
		Example: "ukfast ecloud router create --vpc vpc-abcdef12 --az az-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudRouterCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of router")
	cmd.Flags().String("vpc", "", "ID of VPC")
	cmd.MarkFlagRequired("vpc")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the router has been completely created before continuing on")

	return cmd
}

func ecloudRouterCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateRouterRequest{}
	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}
	createRequest.VPCID, _ = cmd.Flags().GetString("vpc")

	routerID, err := service.CreateRouter(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating router: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(RouterResourceSyncStatusWaitFunc(service, routerID, ecloud.SyncStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for router sync: %s", err)
		}
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
		Short:   "Updates an router",
		Long:    "This command updates one or more routers",
		Example: "ukfast ecloud router update rtr-abcdef12 --name \"my router\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing router")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudRouterUpdate),
	}

	cmd.Flags().String("name", "", "Name of router")

	return cmd
}

func ecloudRouterUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchRouterRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
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
		Short:   "Removes an router",
		Long:    "This command removes one or more routers",
		Example: "ukfast ecloud router delete rtr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing router")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudRouterDelete),
	}
}

func ecloudRouterDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.DeleteRouter(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing router [%s]: %s", arg, err)
		}
	}
	return nil
}

func RouterResourceSyncStatusWaitFunc(service ecloud.ECloudService, routerID string, status ecloud.SyncStatus) helper.WaitFunc {
	return ResourceSyncStatusWaitFunc(func() (ecloud.SyncStatus, error) {
		router, err := service.GetRouter(routerID)
		if err != nil {
			return "", err
		}
		return router.Sync, nil
	}, status)
}
