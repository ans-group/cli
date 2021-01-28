package safedns

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func safednsTemplateRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "sub-commands relating to templates",
	}

	// Child commands
	cmd.AddCommand(safednsTemplateListCmd(f))
	cmd.AddCommand(safednsTemplateShowCmd(f))
	cmd.AddCommand(safednsTemplateCreateCmd(f))
	cmd.AddCommand(safednsTemplateUpdateCmd(f))
	cmd.AddCommand(safednsTemplateDeleteCmd(f))

	// Child root commands
	cmd.AddCommand(safednsTemplateRecordRootCmd(f))

	return cmd
}

func safednsTemplateListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists templates",
		Long:    "This command lists templates",
		Example: "ukfast safedns template list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return safednsTemplateList(c.SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Template name for filtering")

	return cmd
}

func safednsTemplateList(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	templates, err := service.GetTemplates(params)
	if err != nil {
		return fmt.Errorf("Error retrieving templates: %s", err)
	}

	return output.CommandOutput(cmd, OutputSafeDNSTemplatesProvider(templates))

}

func safednsTemplateShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <template: name/id>...",
		Short:   "Shows a template",
		Long:    "This command shows one or more templates",
		Example: "ukfast safedns template show \"main template\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return safednsTemplateShow(c.SafeDNSService(), cmd, args)
		},
	}
}

func safednsTemplateShow(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	var templates []safedns.Template
	for _, arg := range args {
		template, err := getSafeDNSTemplateByNameOrID(service, arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving template [%s]: %s", arg, err)
			continue
		}

		templates = append(templates, template)
	}

	return output.CommandOutput(cmd, OutputSafeDNSTemplatesProvider(templates))
}

func safednsTemplateCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a template",
		Long:    "This command creates a template",
		Example: "ukfast safedns template create --name \"main template\"",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return safednsTemplateCreate(c.SafeDNSService(), cmd, args)
		},
	}
	cmd.Flags().String("name", "", "Name of template")
	cmd.MarkFlagRequired("name")
	cmd.Flags().Bool("default", false, "Specifies template is default")

	return cmd
}

func safednsTemplateCreate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	templateName, _ := cmd.Flags().GetString("name")
	templateDefault, _ := cmd.Flags().GetBool("default")

	createRequest := safedns.CreateTemplateRequest{
		Name:    templateName,
		Default: templateDefault,
	}

	id, err := service.CreateTemplate(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating template: %s", err)
	}

	template, err := service.GetTemplate(id)
	if err != nil {
		return fmt.Errorf("Error retrieving new template: %s", err)
	}

	return output.CommandOutput(cmd, OutputSafeDNSTemplatesProvider([]safedns.Template{template}))
}

func safednsTemplateUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <template: name/id>...",
		Short:   "Updates a template",
		Long:    "This command updates one or more templates",
		Example: "ukfast safedns template update \"main template\" --name \"old template\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return safednsTemplateUpdate(c.SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name of template")
	cmd.Flags().Bool("default", false, "Specifies template is default")

	return cmd
}

func safednsTemplateUpdate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	patchRequest := safedns.PatchTemplateRequest{}

	if cmd.Flags().Changed("name") {
		templateName, _ := cmd.Flags().GetString("name")
		patchRequest.Name = templateName
	}
	if cmd.Flags().Changed("default") {
		templateDefault, _ := cmd.Flags().GetBool("default")
		patchRequest.Default = ptr.Bool(templateDefault)
	}

	var templates []safedns.Template

	for _, arg := range args {
		templateID, err := getSafeDNSTemplateIDByNameOrID(service, arg)
		if err != nil {
			output.OutputWithErrorLevel(err.Error())
			continue
		}

		id, err := service.PatchTemplate(templateID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating template [%s]: %s", arg, err)
			continue
		}

		template, err := service.GetTemplate(id)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated template: %s", err)
			continue
		}

		templates = append(templates, template)
	}

	return output.CommandOutput(cmd, OutputSafeDNSTemplatesProvider(templates))
}

func safednsTemplateDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <template: name/id>...",
		Short:   "Removes a template",
		Long:    "This command removes one or more templates",
		Example: "ukfast safedns template delete \"main template\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			safednsTemplateDelete(c.SafeDNSService(), cmd, args)
			return nil
		},
	}
}

func safednsTemplateDelete(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		templateID, err := getSafeDNSTemplateIDByNameOrID(service, arg)
		if err != nil {
			output.OutputWithErrorLevel(err.Error())
			continue
		}

		err = service.DeleteTemplate(templateID)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing template [%s]: %s", arg, err)
			continue
		}
	}

}
