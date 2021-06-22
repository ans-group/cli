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

func ecloudHostSpecRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hostspec",
		Short: "sub-commands relating to host specs",
	}

	// Child commands
	cmd.AddCommand(ecloudHostSpecListCmd(f))
	cmd.AddCommand(ecloudHostSpecShowCmd(f))

	return cmd
}

func ecloudHostSpecListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists host specs",
		Long:    "This command lists host specs",
		Example: "ukfast ecloud hostspec list",
		RunE:    ecloudCobraRunEFunc(f, ecloudHostSpecList),
	}

	cmd.Flags().String("name", "", "Host spec name for filtering")

	return cmd
}

func ecloudHostSpecList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	specs, err := service.GetHostSpecs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving host specs: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudHostSpecsProvider(specs))
}

func ecloudHostSpecShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <spec: id>...",
		Short:   "Shows an host spec",
		Long:    "This command shows one or more host specs",
		Example: "ukfast ecloud hostspec show hs--abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing host spec")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudHostSpecShow),
	}
}

func ecloudHostSpecShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var specs []ecloud.HostSpec
	for _, arg := range args {
		spec, err := service.GetHostSpec(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving host spec [%s]: %s", arg, err)
			continue
		}

		specs = append(specs, spec)
	}

	return output.CommandOutput(cmd, OutputECloudHostSpecsProvider(specs))
}
