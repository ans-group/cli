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

func ecloudVPNSessionPreSharedKeyRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "presharedkey",
		Short: "sub-commands relating to VPN session pre-shared keys",
	}

	// Child commands
	cmd.AddCommand(ecloudVPNSessionPreSharedKeyShowCmd(f))
	cmd.AddCommand(ecloudVPNSessionPreSharedKeyUpdateCmd(f))

	return cmd
}

func ecloudVPNSessionPreSharedKeyShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Shows VPN session pre-shared keys",
		Long:    "This command shows VPN session pre-shared keys",
		Example: "ans ecloud vpnsession presharedkey show vpns-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN session")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNSessionPreSharedKeyShow),
	}
}

func ecloudVPNSessionPreSharedKeyShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var psks []ecloud.VPNSessionPreSharedKey
	for _, arg := range args {
		psk, err := service.GetVPNSessionPreSharedKey(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving VPN session [%s] pre-shared key: %s", arg, err)
			continue
		}

		psks = append(psks, psk)
	}

	return output.CommandOutput(cmd, OutputECloudVPNSessionPreSharedKeysProvider(psks))
}

func ecloudVPNSessionPreSharedKeyUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <session: id>...",
		Short:   "Updates the pre-shared key for a VPN session",
		Long:    "This command updates the pre-shared key for a VPN session",
		Example: "ans ecloud vpnsession presharedkey update vpns-abcdef12 --psk \"s3curePSK\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN session")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNSessionPreSharedKeyUpdate),
	}

	cmd.Flags().String("psk", "", "Pre-shared key")
	cmd.MarkFlagRequired("psk")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the pre-shared key has been completely updated")

	return cmd
}

func ecloudVPNSessionPreSharedKeyUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	updateRequest := ecloud.UpdateVPNSessionPreSharedKeyRequest{}
	updateRequest.PSK, _ = cmd.Flags().GetString("psk")

	task, err := service.UpdateVPNSessionPreSharedKey(args[0], updateRequest)
	if err != nil {
		return fmt.Errorf("Error updating VPN session pre-shared key: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for task to complete for VPN session: %s", err)
		}
	}

	return nil
}
