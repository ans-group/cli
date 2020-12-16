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

func ecloudInstanceCredentialRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credential",
		Short: "sub-commands relating to instance credentials",
	}

	// Child commands
	cmd.AddCommand(ecloudInstanceCredentialListCmd(f))

	return cmd
}

func ecloudInstanceCredentialListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists instance credentials",
		Long:    "This command lists instance credentials",
		Example: "ukfast ecloud instance credential list i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceCredentialList),
	}

	cmd.Flags().String("name", "", "Credential name for filtering")

	return cmd
}

func ecloudInstanceCredentialList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	helper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, helper.NewStringFilterFlag("name", "name"))

	credentials, err := service.GetInstanceCredentials(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving instance credentials: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudCredentialsProvider(credentials))
}
