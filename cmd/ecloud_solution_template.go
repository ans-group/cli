package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSolutionTemplateRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "sub-commands relating to solution templates",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionTemplateListCmd())
	cmd.AddCommand(ecloudSolutionTemplateShowCmd())

	return cmd
}

func ecloudSolutionTemplateListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists solution templates",
		Long:    "This command lists solution templates",
		Example: "ukfast ecloud solution template list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionTemplateList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionTemplateList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
	}

	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	templates, err := service.GetSolutionTemplates(solutionID, params)
	if err != nil {
		output.Fatalf("Error retrieving solution templates: %s", err)
		return
	}

	outputECloudTemplates(templates)
}

func ecloudSolutionTemplateShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <solution: id> <template: name>...",
		Short:   "Shows a solution template",
		Long:    "This command shows one or more solution templates",
		Example: "ukfast ecloud solution template show 123 foo",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}
			if len(args) < 2 {
				return errors.New("Missing template")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionTemplateShow(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionTemplateShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
	}

	var templates []ecloud.Template

	for _, arg := range args[1:] {
		template, err := service.GetSolutionTemplate(solutionID, arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving solution template [%s]: %s", arg, err)
			continue
		}

		templates = append(templates, template)
	}

	outputECloudTemplates(templates)
}
