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

func ecloudFloatingIPRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "floatingip",
		Short: "sub-commands relating to floating IPs",
	}

	// Child commands
	cmd.AddCommand(ecloudFloatingIPListCmd(f))
	cmd.AddCommand(ecloudFloatingIPShowCmd(f))
	cmd.AddCommand(ecloudFloatingIPCreateCmd(f))
	cmd.AddCommand(ecloudFloatingIPUpdateCmd(f))
	cmd.AddCommand(ecloudFloatingIPDeleteCmd(f))
	cmd.AddCommand(ecloudFloatingIPAssignCmd(f))
	cmd.AddCommand(ecloudFloatingIPUnassignCmd(f))

	return cmd
}

func ecloudFloatingIPListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists floating IPs",
		Long:    "This command lists floating IPs",
		Example: "ukfast ecloud floatingip list",
		RunE:    ecloudCobraRunEFunc(f, ecloudFloatingIPList),
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
		RunE: ecloudCobraRunEFunc(f, ecloudFloatingIPShow),
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

func ecloudFloatingIPCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a floating IP",
		Long:    "This command creates a floating IP address",
		Example: "ukfast ecloud floatingip create --vpc vpc-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudFloatingIPCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of floating IP")
	cmd.Flags().String("vpc", "", "ID of VPC")
	cmd.MarkFlagRequired("vpc")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the floating IP has been completely created")

	return cmd
}

func ecloudFloatingIPCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateFloatingIPRequest{}
	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}
	createRequest.VPCID, _ = cmd.Flags().GetString("vpc")

	fipID, err := service.CreateFloatingIP(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating floating IP: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(FloatingIPResourceSyncStatusWaitFunc(service, fipID, ecloud.SyncStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for floating IP sync: %s", err)
		}
	}

	fip, err := service.GetFloatingIP(fipID)
	if err != nil {
		return fmt.Errorf("Error retrieving new floating IP: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFloatingIPsProvider([]ecloud.FloatingIP{fip}))
}

func ecloudFloatingIPUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <fip: id>...",
		Short:   "Updates a floating IP",
		Long:    "This command updates one or more floating IPs",
		Example: "ukfast ecloud floatingip update fip-abcdef12 --name \"my fip\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing floating IP")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFloatingIPUpdate),
	}

	cmd.Flags().String("name", "", "Name of floating IP")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the floating IP has been completely updated")

	return cmd
}

func ecloudFloatingIPUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchFloatingIPRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var fips []ecloud.FloatingIP
	for _, arg := range args {
		err := service.PatchFloatingIP(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating floating IP [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(FloatingIPResourceSyncStatusWaitFunc(service, arg, ecloud.SyncStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for floating IP [%s] sync: %s", arg, err)
				continue
			}
		}

		fip, err := service.GetFloatingIP(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated floating IP [%s]: %s", arg, err)
			continue
		}

		fips = append(fips, fip)
	}

	return output.CommandOutput(cmd, OutputECloudFloatingIPsProvider(fips))
}

func ecloudFloatingIPDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <fip: id>...",
		Short:   "Removes a floating IP",
		Long:    "This command removes one or more floating IPs",
		Example: "ukfast ecloud floatingip delete fip-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing floating IP")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFloatingIPDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the floating IP has been completely removed")

	return cmd
}

func ecloudFloatingIPDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.DeleteFloatingIP(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing floating IP [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(FloatingIPNotFoundWaitFunc(service, arg))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for removal of floating IP [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}

func ecloudFloatingIPAssignCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "assign <fip: id>",
		Short:   "Assigns a floating IP to a resource",
		Long:    "This command assigns a floating IP to a resource",
		Example: "ukfast ecloud floatingip assign fip-abcdef12 --resource i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing floating IP")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFloatingIPAssign),
	}

	cmd.Flags().String("resource", "", "ID of resource to assign")
	cmd.MarkFlagRequired("resource")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the floating IP has been completely assigned")

	return cmd
}

func ecloudFloatingIPAssign(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	fipID := args[0]
	resource, _ := cmd.Flags().GetString("resource")
	req := ecloud.AssignFloatingIPRequest{
		ResourceID: resource,
	}

	err := service.AssignFloatingIP(fipID, req)
	if err != nil {
		return fmt.Errorf("Error assigning floating IP to resource: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(FloatingIPResourceSyncStatusWaitFunc(service, fipID, ecloud.SyncStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for floating IP sync: %s", err)
		}
	}

	fip, err := service.GetFloatingIP(fipID)
	if err != nil {
		return fmt.Errorf("Error retrieving new floating IP: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFloatingIPsProvider([]ecloud.FloatingIP{fip}))
}

func ecloudFloatingIPUnassignCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unassign <fip: id>...",
		Short:   "Unassigns a floating IP",
		Long:    "This command unassigns one or more floating IPs from connected resources",
		Example: "ukfast ecloud floatingip unassign fip-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing floating IP")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFloatingIPUnassign),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the floating IP has been completely unassigned")

	return cmd
}

func ecloudFloatingIPUnassign(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.UnassignFloatingIP(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error unassigning floating IP [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(FloatingIPResourceSyncStatusWaitFunc(service, arg, ecloud.SyncStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for floating IP [%s] sync: %s", arg, err)
				continue
			}
		}
	}
	return nil
}

func FloatingIPResourceSyncStatusWaitFunc(service ecloud.ECloudService, fipID string, status ecloud.SyncStatus) helper.WaitFunc {
	return ResourceSyncStatusWaitFunc(func() (ecloud.SyncStatus, error) {
		fip, err := service.GetFloatingIP(fipID)
		if err != nil {
			return "", err
		}
		return fip.Sync.Status, nil
	}, status)
}

func FloatingIPNotFoundWaitFunc(service ecloud.ECloudService, fipID string) helper.WaitFunc {
	return func() (finished bool, err error) {
		_, err = service.GetFloatingIP(fipID)
		if err != nil {
			switch err.(type) {
			case *ecloud.FloatingIPNotFoundError:
				return true, nil
			default:
				return false, fmt.Errorf("Failed to retrieve floating IP [%s]: %s", fipID, err)
			}
		}

		return false, nil
	}
}
