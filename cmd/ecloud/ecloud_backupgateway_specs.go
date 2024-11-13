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

func ecloudBackupGatewaySpecificationRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spec",
		Short: "sub-commands relating to backup gateway specifications",
	}

	// Child commands
	cmd.AddCommand(ecloudBackupGatewaySpecificationListCmd(f))
	cmd.AddCommand(ecloudBackupGatewaySpecificationShowCmd(f))

	return cmd
}

func ecloudBackupGatewaySpecificationListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists Backup gateway specifications",
		Example: "ans ecloud backupgateway spec list",
		RunE:    ecloudCobraRunEFunc(f, ecloudBackupGatewaySpecificationList),
	}

	cmd.Flags().String("name", "", "Backup gateway specification name for filtering")

	return cmd
}

func ecloudBackupGatewaySpecificationList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	specs, err := service.GetBackupGatewaySpecifications(params)
	if err != nil {
		return fmt.Errorf("Error retrieving backup gateway specifications: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudBackupGatewaySpecificationsProvider(specs))
}

func ecloudBackupGatewaySpecificationShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <specification: id>...",
		Short:   "Show details of a backup gateway specification",
		Example: "ans ecloud backupgateway spec show bgws-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing backup gateway specification")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudBackupGatewaySpecificationShow),
	}
}

func ecloudBackupGatewaySpecificationShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var specs []ecloud.BackupGatewaySpecification
	for _, arg := range args {
		spec, err := service.GetBackupGatewaySpecification(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving backup gateway specification [%s]: %s", arg, err)
			continue
		}

		specs = append(specs, spec)
	}

	return output.CommandOutput(cmd, OutputECloudBackupGatewaySpecificationsProvider(specs))
}
