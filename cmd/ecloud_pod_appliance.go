package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudPodApplianceRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "appliance",
		Short: "sub-commands relating to pod appliances",
	}

	// Child commands
	cmd.AddCommand(ecloudPodApplianceListCmd())

	return cmd
}

func ecloudPodApplianceListCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudPodApplianceList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudPodApplianceList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	podID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid pod ID [%s]", args[0])
		return
	}

	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	appliances, err := service.GetPodAppliances(podID, params)
	if err != nil {
		output.Fatalf("Error retrieving pod appliances: %s", err)
		return
	}

	outputECloudAppliances(appliances)
}
