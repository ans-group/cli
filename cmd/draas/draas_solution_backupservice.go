package draas

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/draas"
	"github.com/spf13/cobra"
)

func draasSolutionBackupServiceRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backupservice",
		Short: "sub-commands relating to solution backup services",
	}

	// Child commands
	cmd.AddCommand(draasSolutionBackupServiceShowCmd(f))
	cmd.AddCommand(draasSolutionBackupServiceResetCredentialsCmd(f))

	return cmd
}

func draasSolutionBackupServiceShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <solution: id>",
		Short:   "Shows the backup service for a solution",
		Long:    "This command shows the backup service for a solution",
		Example: "ans draas solution backupservice show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return draasSolutionBackupServiceShow(c.DRaaSService(), cmd, args)
		},
	}
}

func draasSolutionBackupServiceShow(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	backupService, err := service.GetSolutionBackupService(args[0])
	if err != nil {
		return fmt.Errorf("Error retrieving solution backup service: %s", err)
	}

	return output.CommandOutput(cmd, BackupServiceCollection([]draas.BackupService{backupService}))
}

func draasSolutionBackupServiceResetCredentialsCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resetcredentials <solution: id>",
		Short:   "Resets the backup service credentials for a solution",
		Long:    "This command resets the backup service credentials for a solution",
		Example: "ans draas solution backupservice resetcredentials 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return draasSolutionBackupServiceResetCredentials(c.DRaaSService(), cmd, args)
		},
	}

	cmd.Flags().String("password", "", "New password to set")
	cmd.MarkFlagRequired("password")

	return cmd
}

func draasSolutionBackupServiceResetCredentials(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	password, _ := cmd.Flags().GetString("password")

	req := draas.ResetBackupServiceCredentialsRequest{
		Password: password,
	}

	err := service.ResetSolutionBackupServiceCredentials(args[0], req)
	if err != nil {
		return fmt.Errorf("Error resetting credentials for solution backup service: %s", err)
	}

	return nil
}
