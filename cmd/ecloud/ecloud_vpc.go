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

func ecloudVPCRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vpc",
		Short: "sub-commands relating to VPCs",
	}

	// Child commands
	cmd.AddCommand(ecloudVPCListCmd(f))
	cmd.AddCommand(ecloudVPCShowCmd(f))
	cmd.AddCommand(ecloudVPCCreateCmd(f))
	cmd.AddCommand(ecloudVPCUpdateCmd(f))
	cmd.AddCommand(ecloudVPCDeleteCmd(f))
	cmd.AddCommand(ecloudVPCDeployDefaultsCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudVPCVolumeRootCmd(f))
	cmd.AddCommand(ecloudVPCInstanceRootCmd(f))
	cmd.AddCommand(ecloudVPCTaskRootCmd(f))

	return cmd
}

func ecloudVPCListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists VPCs",
		Long:    "This command lists VPCs",
		Example: "ans ecloud vpc list",
		RunE:    ecloudCobraRunEFunc(f, ecloudVPCList),
	}

	cmd.Flags().String("name", "", "VPC name for filtering")

	return cmd
}

func ecloudVPCList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	vpcs, err := service.GetVPCs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving VPCs: %s", err)
	}

	return output.CommandOutput(cmd, VPCCollection(vpcs))
}

func ecloudVPCShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <vpc: id>...",
		Short:   "Shows a VPC",
		Long:    "This command shows one or more VPCs",
		Example: "ans ecloud vpc show vpc-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing vpc")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPCShow),
	}
}

func ecloudVPCShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var vpcs []ecloud.VPC
	for _, arg := range args {
		vpc, err := service.GetVPC(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving VPC [%s]: %s", arg, err)
			continue
		}

		vpcs = append(vpcs, vpc)
	}

	return output.CommandOutput(cmd, VPCCollection(vpcs))
}

func ecloudVPCCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a VPC",
		Long:    "This command creates a VPC",
		Example: "ans ecloud vpc create",
		RunE:    ecloudCobraRunEFunc(f, ecloudVPCCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of VPC")
	cmd.Flags().String("region", "", "ID of region")
	cmd.MarkFlagRequired("region")
	cmd.Flags().Int("client-id", 0, "ID of client")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPC has been completely created")

	return cmd
}

func ecloudVPCCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {

	name, _ := cmd.Flags().GetString("name")
	regionID, _ := cmd.Flags().GetString("region")
	clientID, _ := cmd.Flags().GetInt("client-id")

	createRequest := ecloud.CreateVPCRequest{
		Name:     name,
		RegionID: regionID,
		ClientID: clientID,
	}

	vpcID, err := service.CreateVPC(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating VPC: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(VPCResourceSyncStatusWaitFunc(service, vpcID, ecloud.SyncStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for VPC sync: %s", err)
		}
	}

	vpc, err := service.GetVPC(vpcID)
	if err != nil {
		return fmt.Errorf("Error retrieving new VPC: %s", err)
	}

	return output.CommandOutput(cmd, VPCCollection([]ecloud.VPC{vpc}))
}

func ecloudVPCUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <vpc: id>...",
		Short:   "Updates a VPC",
		Long:    "This command updates one or more VPCs",
		Example: "ans ecloud vpc update vpc-abcdef12 --name \"my vpc\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing vpc")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPCUpdate),
	}

	cmd.Flags().String("name", "", "Name of VPC")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPC has been completely updated")

	return cmd
}

func ecloudVPCUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchVPCRequest{}

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		patchRequest.Name = name
	}

	var vpcs []ecloud.VPC
	for _, arg := range args {
		err := service.PatchVPC(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating VPC [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(VPCResourceSyncStatusWaitFunc(service, arg, ecloud.SyncStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for VPC [%s] sync: %s", arg, err)
				continue
			}
		}

		vpc, err := service.GetVPC(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated VPC [%s]: %s", arg, err)
			continue
		}

		vpcs = append(vpcs, vpc)
	}

	return output.CommandOutput(cmd, VPCCollection(vpcs))
}

func ecloudVPCDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <vpc: id>...",
		Short:   "Removes a VPC",
		Long:    "This command removes one or more VPCs",
		Example: "ans ecloud vpc delete vpc-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing vpc")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			ecloudVPCDelete(c.ECloudService(), cmd, args)
			return nil
		},
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPC has been completely removed")
	cmd.Flags().Bool("recursive", false, "Recursively delete all resources within the VPC before deleting the VPC itself")
	cmd.Flags().Bool("dangerously-recursively-delete", false, "Skip interactive confirmation when using recursive deletion (for automation)")

	return cmd
}

func ecloudVPCDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	recursive, _ := cmd.Flags().GetBool("recursive")
	dangerousFlag, _ := cmd.Flags().GetBool("dangerously-recursively-delete")
	waitFlag, _ := cmd.Flags().GetBool("wait")

	for _, vpcID := range args {
		if recursive {
			if !dangerousFlag {
				confirmed, err := confirmVPCRecursiveDeletion(vpcID)
				if err != nil {
					output.OutputWithErrorLevelf("Error getting confirmation for VPC [%s]: %s", vpcID, err)
					continue
				}
				if !confirmed {
					fmt.Printf("VPC [%s] deletion cancelled\n", vpcID)
					continue
				}
			}

			err := deleteVPCResourcesRecursively(service, vpcID)
			if err != nil {
				output.OutputWithErrorLevelf("ecloud: Error deleting resources for VPC [%s]: %s", vpcID, err)
				continue
			}
		} else {
			err := service.DeleteVPC(vpcID)
			if err != nil {
				output.OutputWithErrorLevelf("ecloud: Error removing VPC [%s]: %s", vpcID, err)
				continue
			}

			if waitFlag {
				err := helper.WaitForCommand(VPCNotFoundWaitFunc(service, vpcID))
				if err != nil {
					output.OutputWithErrorLevelf("ecloud: Error waiting for removal of VPC [%s]: %s", vpcID, err)
					continue
				}
			}
		}
	}
}

func ecloudVPCDeployDefaultsCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "deploydefaults <vpc: id>...",
		Short:   "Deploys default resources for a VPC",
		Long:    "This command deploys default resources for one or more VPCs",
		Example: "ans ecloud vpc deploydefaults vpc-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing vpc")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPCDeployDefaults),
	}
}

func ecloudVPCDeployDefaults(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.DeployVPCDefaults(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error deploying default resources for VPC [%s]: %s", arg, err)
			continue
		}
	}

	return nil
}

func VPCResourceSyncStatusWaitFunc(service ecloud.ECloudService, vpcID string, status ecloud.SyncStatus) helper.WaitFunc {
	return ResourceSyncStatusWaitFunc(func() (ecloud.SyncStatus, error) {
		vpc, err := service.GetVPC(vpcID)
		if err != nil {
			return "", err
		}
		return vpc.Sync.Status, nil
	}, status)
}

func VPCNotFoundWaitFunc(service ecloud.ECloudService, vpcID string) helper.WaitFunc {
	return func() (finished bool, err error) {
		_, err = service.GetVPC(vpcID)
		if err != nil {
			switch err.(type) {
			case *ecloud.VPCNotFoundError:
				return true, nil
			default:
				return false, fmt.Errorf("Failed to retrieve VPC [%s]: %s", vpcID, err)
			}
		}

		return false, nil
	}
}
