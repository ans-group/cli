package draas

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	flaghelper "github.com/ukfast/cli/internal/pkg/helper/flag"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/draas"
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
		Example: "ukfast draas solution backupresource list 00000000-0000-0000-0000-000000000000",
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

			return draasSolutionBackupResourceList(c.DRaaSService(), cmd, args)
		},
	}
}

func draasSolutionBackupResourceList(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	params, err := flaghelper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	resources, err := service.GetSolutionBackupResources(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving solution backup resources: %s", err)
	}

	return output.CommandOutput(cmd, OutputDRaaSBackupResourcesProvider(resources))
}
