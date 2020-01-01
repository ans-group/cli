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
	cmd.AddCommand(ecloudPodTemplateShowCmd())
	cmd.AddCommand(ecloudPodTemplateUpdateCmd())
	cmd.AddCommand(ecloudPodTemplateDeleteCmd())

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

func ecloudPodTemplateShowCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudPodTemplateShow(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudPodTemplateShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	podID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid pod ID [%s]", args[0])
		return
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

	outputECloudTemplates(templates)
}

func ecloudPodTemplateUpdateCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudPodTemplateUpdate(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name for template")

	return cmd
}

func ecloudPodTemplateUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	podID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid pod ID [%s]", args[0])
		return
	}

	templateName := args[1]

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		patchRequest := ecloud.RenameTemplateRequest{
			Destination: name,
		}

		err = service.RenamePodTemplate(podID, templateName, patchRequest)
		if err != nil {
			output.Fatalf("Error updating pod template: %s", err)
			return
		}

		err := WaitForCommand(PodTemplateExistsWaitFunc(service, podID, name, true))
		if err != nil {
			output.Fatalf("Error waiting for pod template update: %s", err)
			return
		}

		templateName = name
	}

	template, err := service.GetPodTemplate(podID, templateName)
	if err != nil {
		output.Fatalf("Error retrieving updated pod template: %s", err)
		return
	}

	outputECloudTemplates([]ecloud.Template{template})
}

func ecloudPodTemplateDeleteCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudPodTemplateDelete(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the template has been completely created before continuing on")

	return cmd
}

func ecloudPodTemplateDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	podID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid pod ID [%s]", args[0])
		return
	}

	for _, arg := range args[1:] {
		err = service.DeletePodTemplate(podID, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing pod template [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := WaitForCommand(PodTemplateExistsWaitFunc(service, podID, arg, false))
			if err != nil {
				output.OutputWithErrorLevelf("Error removing pod template [%s]: %s", arg, err)
			}
		}
	}
}
