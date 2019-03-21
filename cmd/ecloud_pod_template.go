package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudPodTemplateRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "sub-commands relating to pod templates",
	}

	// Child commands
	cmd.AddCommand(ecloudPodTemplateListCmd())

	return cmd
}

func ecloudPodTemplateListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists pod templates",
		Long:    "This command lists pod templates",
		Example: "ukfast ecloud pod template list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing pod")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudPodTemplateList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudPodTemplateList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	podID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid pod ID [%s]", args[0])
		return
	}

	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	templates, err := service.GetPodTemplates(podID, params)
	if err != nil {
		output.Fatalf("Error retrieving pod templates: %s", err)
		return
	}

	outputECloudTemplates(templates)
}
