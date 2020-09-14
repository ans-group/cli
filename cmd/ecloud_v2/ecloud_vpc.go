package ecloud_v2

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudVPCRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vpc",
		Short: "sub-commands relating to vpcs",
	}

	// Child commands
	cmd.AddCommand(ecloudVPCListCmd(f))
	cmd.AddCommand(ecloudVPCShowCmd(f))

	return cmd
}

func ecloudVPCListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists vpcs",
		Long:    "This command lists vpcs",
		Example: "ukfast ecloud vpc list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudVPCList(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "VPC name for filtering")

	return cmd
}

func ecloudVPCList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	helper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, "name", "name")

	vpcs, err := service.GetVPCs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving vpcs: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVPCsProvider(vpcs))
}

func ecloudVPCShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <vpc: id>...",
		Short:   "Shows a vpc",
		Long:    "This command shows one or more vpcs",
		Example: "ukfast ecloud vpc show 123",
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

			return ecloudVPCShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudVPCShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var vpcs []ecloud.VPC
	for _, arg := range args {
		vpc, err := service.GetVPC(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving vpc [%s]: %s", arg, err)
			continue
		}

		vpcs = append(vpcs, vpc)
	}

	return output.CommandOutput(cmd, OutputECloudVPCsProvider(vpcs))
}
