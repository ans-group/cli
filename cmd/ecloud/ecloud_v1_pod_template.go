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

func ecloudPodTemplateRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "sub-commands relating to pod templates",
	}

	// Child commands
	cmd.AddCommand(ecloudPodTemplateListCmd(f))
	cmd.AddCommand(ecloudPodTemplateShowCmd(f))
	cmd.AddCommand(ecloudPodTemplateUpdateCmd(f))
	cmd.AddCommand(ecloudPodTemplateDeleteCmd(f))

	return cmd
}

func ecloudPodTemplateListCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudPodTemplateList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudPodTemplateList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	podID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid pod ID [%s]", args[0])
	}

	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	templates, err := service.GetPodTemplates(podID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving pod templates: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudTemplatesProvider(templates))
}

func ecloudPodTemplateShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <pod: id> <template: name>...",
		Short:   "Shows a pod template",
		Long:    "This command shows one or more pod templates",
		Example: "ukfast ecloud pod template show 123 foo",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing pod")
			}
			if len(args) < 2 {
				return errors.New("Missing template")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudPodTemplateShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudPodTemplateShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	podID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid pod ID [%s]", args[0])
	}

	var templates []ecloud.Template

	for _, arg := range args[1:] {
		template, err := service.GetPodTemplate(podID, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving pod template [%s]: %s", arg, err)
			continue
		}

		templates = append(templates, template)
	}

	return output.CommandOutput(cmd, OutputECloudTemplatesProvider(templates))
}

func ecloudPodTemplateUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <pod: id> <template: name>...",
		Short:   "Updates a pod template",
		Long:    "This command updates a pod template",
		Example: "ukfast ecloud pod template update 123 foo --name \"bar\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing pod")
			}
			if len(args) < 2 {
				return errors.New("Missing template")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudPodTemplateUpdate(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name for template")

	return cmd
}

func ecloudPodTemplateUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	podID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid pod ID [%s]", args[0])
	}

	templateName := args[1]

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		patchRequest := ecloud.RenameTemplateRequest{
			Destination: name,
		}

		err = service.RenamePodTemplate(podID, templateName, patchRequest)
		if err != nil {
			return fmt.Errorf("Error updating pod template: %s", err)
		}

		err := helper.WaitForCommand(PodTemplateExistsWaitFunc(service, podID, name, true))
		if err != nil {
			return fmt.Errorf("Error waiting for pod template update: %s", err)
		}

		templateName = name
	}

	template, err := service.GetPodTemplate(podID, templateName)
	if err != nil {
		return fmt.Errorf("Error retrieving updated pod template: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudTemplatesProvider([]ecloud.Template{template}))
}

func ecloudPodTemplateDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <pod: id> <template: name>...",
		Short:   "Removes a pod template ",
		Long:    "This command removes one or more pod templates",
		Example: "ukfast ecloud pod template delete 123 foo",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing pod")
			}
			if len(args) < 2 {
				return errors.New("Missing template")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudPodTemplateDelete(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the template has been completely created before continuing on")

	return cmd
}

func ecloudPodTemplateDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	podID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid pod ID [%s]", args[0])
	}

	for _, arg := range args[1:] {
		err = service.DeletePodTemplate(podID, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing pod template [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(PodTemplateExistsWaitFunc(service, podID, arg, false))
			if err != nil {
				output.OutputWithErrorLevelf("Error removing pod template [%s]: %s", arg, err)
			}
		}
	}

	return nil
}
