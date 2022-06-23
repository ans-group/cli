package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

func ecloudPodApplianceRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "appliance",
		Short: "sub-commands relating to pod appliances",
	}

	// Child commands
	cmd.AddCommand(ecloudPodApplianceListCmd(f))

	return cmd
}

func ecloudPodApplianceListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists pod appliances",
		Long:    "This command lists pod appliances",
		Example: "ukfast ecloud pod appliance list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing pod")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudPodApplianceList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudPodApplianceList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	podID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid pod ID [%s]", args[0])
	}

	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	appliances, err := service.GetPodAppliances(podID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving pod appliances: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudAppliancesProvider(appliances))
}
