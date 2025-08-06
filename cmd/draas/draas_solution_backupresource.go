package draas

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/draas"
	"github.com/spf13/cobra"
)

func draasSolutionBackupResourceRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backupresource",
		Short: "sub-commands relating to solution backup resources",
	}

	// Child commands
	cmd.AddCommand(draasSolutionBackupResourceListCmd(f))

	// Child root commands

	return cmd
}

func draasSolutionBackupResourceListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <solution: id>",
		Short:   "Lists a solution",
		Long:    "This command lists the backup resources for a solution",
		Example: "ans draas solution backupresource list 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing solution")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return draasSolutionBackupResourceList(c.DRaaSService(), cmd, args)
		},
	}
}

func draasSolutionBackupResourceList(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	resources, err := service.GetSolutionBackupResources(args[0], params)
	if err != nil {
		return fmt.Errorf("error retrieving solution backup resources: %s", err)
	}

	return output.CommandOutput(cmd, BackupResourceCollection(resources))
}
