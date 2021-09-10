package ecloud

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudVPNSessionPreSharedKeyRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "presharedkey",
		Short: "sub-commands relating to VPN session pre-shared keys",
	}

	// Child commands
	cmd.AddCommand(ecloudVPNSessionPreSharedKeyShowCmd(f))

	return cmd
}

func ecloudVPNSessionPreSharedKeyShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Shows VPN session pre-shared keys",
		Long:    "This command shows VPN session pre-shared keys",
		Example: "ukfast ecloud vpnsession presharedkey show vpns-abcdef12",
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
