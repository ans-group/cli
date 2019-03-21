package cmd

import (
	"errors"

	"github.com/ukfast/sdk-go/pkg/ptr"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudTemplateRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "sub-commands relating to templates",
	}

	// Child commands
	cmd.AddCommand(ecloudTemplateListCmd())
	cmd.AddCommand(ecloudTemplateShowCmd())
	cmd.AddCommand(ecloudTemplateUpdateCmd())
	cmd.AddCommand(ecloudTemplateDeleteCmd())

	// Child root commands
	cmd.AddCommand(ecloudTemplateDiskRootCmd())

	return cmd
}

func ecloudTemplateListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists templates",
		Long:    "This command lists templates",
		Example: "ukfast ecloud template list",
		Run: func(cmd *cobra.Command, args []string) {
			ecloudTemplateList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudTemplateList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	templates, err := service.GetTemplates(params)
	if err != nil {
		output.Fatalf("Error retrieving templates: %s", err)
		return
	}

	outputECloudTemplates(templates)
}

func ecloudTemplateShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <template: name>...",
		Short:   "Shows a template",
		Long:    "This command shows one or more templates",
		Example: "ukfast ecloud template \"test template\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudTemplateShow(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudTemplateShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	var templates []ecloud.Template
	for _, arg := range args {
		template, err := service.GetTemplate(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving template [%s]: %s", arg, err)
			continue
		}

		templates = append(templates, template)
	}

	outputECloudTemplates(templates)
}

func ecloudTemplateUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <template: name>",
		Short:   "Updates a template",
		Long:    "This command updates a template",
		Example: "ukfast ecloud template update \"current name\" --name \"new name\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudTemplateUpdate(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name for template")
	cmd.MarkFlagRequired("name")
	cmd.Flags().Int("solution", 0, "Solution ID for template")

	return cmd
}

func ecloudTemplateUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	renameRequest := ecloud.RenameTemplateRequest{}

	newTemplateName, _ := cmd.Flags().GetString("name")
	renameRequest.NewTemplateName = newTemplateName

	if cmd.Flags().Changed("solution") {
		solutionID, _ := cmd.Flags().GetInt("solution")
		renameRequest.SolutionID = ptr.Int(solutionID)
	}

	err := service.RenameTemplate(args[0], renameRequest)
	if err != nil {
		output.Fatalf("Error updating template [%s]: %s", args[0], err)
		return
	}

	template, err := service.GetTemplate(newTemplateName)
	if err != nil {
		output.Fatalf("Error retrieving updated template: %s", err)
		return
	}

	outputECloudTemplates([]ecloud.Template{template})
}

func ecloudTemplateDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <template: name>...",
		Short:   "Removes a template",
		Long:    "This command removes one or more templates",
		Example: "ukfast ecloud template delete ukfast.co.uk\nukfast ecloud template delete 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudTemplateDelete(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudTemplateDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		err := service.DeleteTemplate(arg)
		if err != nil {
			OutputWithErrorLevelf("Error removing template [%s]: %s", arg, err)
		}
	}
}
