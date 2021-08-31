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

func ecloudVPNSessionCredentialRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credential",
		Short: "sub-commands relating to VPN session credentials",
	}

	// Child commands
	cmd.AddCommand(ecloudVPNSessionCredentialListCmd(f))

	return cmd
}

func ecloudVPNSessionCredentialListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists VPN session credentials",
		Long:    "This command lists VPN session credentials",
		Example: "ukfast ecloud vpnsession credential list vpns-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN session")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNSessionCredentialList),
	}

	cmd.Flags().String("name", "", "Credential name for filtering")

	return cmd
}

func ecloudVPNSessionCredentialList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	credentials, err := service.GetVPNSessionCredentials(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving VPN session credentials: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudCredentialsProvider(credentials))
}
