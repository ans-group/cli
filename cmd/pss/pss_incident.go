package pss

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssIncidentRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "incident",
		Short: "sub-commands relating to incident cases",
	}

	// Child commands
	cmd.AddCommand(pssIncidentListCmd(f))
	cmd.AddCommand(pssIncidentShowCmd(f))
	cmd.AddCommand(pssIncidentCreateCmd(f))
	cmd.AddCommand(pssIncidentCloseCmd(f))

	// Child root commands
	cmd.AddCommand(pssIncidentTypeRootCmd(f))
	cmd.AddCommand(pssIncidentImpactRootCmd(f))

	// Additional root commands (generic case commands)
	cmd.AddCommand(pssCaseUpdateRootCmd(f))
	cmd.AddCommand(pssCaseCategoryRootCmd(f))

	return cmd
}

func pssIncidentListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists incident cases",
		Long:    "This command lists incident cases",
		Example: "ans pss incident list",
		RunE:    pssCobraRunEFunc(f, pssIncidentList),
	}
}

func pssIncidentList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	incidents, err := service.GetIncidentCases(params)
	if err != nil {
		return err
	}

	return output.CommandOutput(cmd, OutputPSSIncidentCasesProvider(incidents))
}

func pssIncidentShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <incident: id>...",
		Short:   "Shows an incident",
		Long:    "This command shows one or more incidents",
		Example: "ans pss incident show INC123456",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing incident")
			}

			return nil
		},
		RunE: pssCobraRunEFunc(f, pssIncidentShow),
	}
}

func pssIncidentShow(service pss.PSSService, cmd *cobra.Command, args []string) error {
	var incidents []pss.IncidentCase
	for _, arg := range args {
		incident, err := service.GetIncidentCase(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving incident [%s]: %s", arg, err)
			continue
		}

		incidents = append(incidents, incident)
	}

	return output.CommandOutput(cmd, OutputPSSIncidentCasesProvider(incidents))
}

func pssIncidentCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an incident",
		Long:    "This command creates a new incident",
		Example: "ans pss incident create --title 'test incident' --description 'test incident' --type Fault --category 70b67a49-7ace-4146-a295-11590a0b1203 --supported-service 1b684067-a587-4997-acb9-da2f4cb7be81 --impact Minor",
		RunE:    pssCobraRunEFunc(f, pssIncidentCreate),
	}

	// Setup flags
	cmd.Flags().String("title", "", "Specifies the title for incident case")
	cmd.MarkFlagRequired("title")
	cmd.Flags().String("description", "", "Specifies the description for incident case")
	cmd.MarkFlagRequired("description")
	cmd.Flags().String("type", "", "Specifies the type of incident case")
	cmd.MarkFlagRequired("type")
	cmd.Flags().String("category", "", "Category ID for incident case")
	cmd.MarkFlagRequired("category")
	cmd.Flags().String("supported-service", "", "Supported service ID for incident case")
	cmd.MarkFlagRequired("supported-service")
	cmd.Flags().String("impact", "", "Impact for incident case")
	cmd.MarkFlagRequired("impact")
	cmd.Flags().Bool("security", false, "Specifies whether incident case is a security incident")
	cmd.Flags().String("customer-reference", "", "Specifies the customer reference for incident case")
	cmd.Flags().Int("contact", 0, "Contact ID for incident case")

	return cmd
}

func pssIncidentCreate(service pss.PSSService, cmd *cobra.Command, args []string) error {
	createIncidentCase := pss.CreateIncidentCaseRequest{}

	createIncidentCase.Title, _ = cmd.Flags().GetString("title")
	createIncidentCase.Description, _ = cmd.Flags().GetString("description")
	createIncidentCase.IsSecurity, _ = cmd.Flags().GetBool("security")
	createIncidentCase.CustomerReference, _ = cmd.Flags().GetString("customer-reference")
	createIncidentCase.ContactID, _ = cmd.Flags().GetInt("contact")
	createIncidentCase.CategoryID, _ = cmd.Flags().GetString("category")
	createIncidentCase.SupportedServiceID, _ = cmd.Flags().GetString("supported-service")

	incidentCaseType, _ := cmd.Flags().GetString("type")
	parsedIncidentCaseType, err := pss.IncidentCaseTypeEnum.Parse(incidentCaseType)
	if err != nil {
		return err
	}
	createIncidentCase.Type = parsedIncidentCaseType

	incidentCaseImpact, _ := cmd.Flags().GetString("impact")
	parsedIncidentCaseImpact, err := pss.IncidentCaseImpactEnum.Parse(incidentCaseImpact)
	if err != nil {
		return err
	}
	createIncidentCase.Impact = parsedIncidentCaseImpact

	incidentID, err := service.CreateIncidentCase(createIncidentCase)
	if err != nil {
		return fmt.Errorf("Error creating incident: %s", err)
	}

	incident, err := service.GetIncidentCase(incidentID)
	if err != nil {
		return fmt.Errorf("Error retrieving new incident: %s", err)
	}

	return output.CommandOutput(cmd, OutputPSSIncidentCasesProvider([]pss.IncidentCase{incident}))
}

func pssIncidentCloseCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "close <incident: id>...",
		Short:   "Closes an incident",
		Long:    "This command closes one or more incidents",
		Example: "ans pss incident close CHG123456",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing incident")
			}

			return nil
		},
		RunE: pssCobraRunEFunc(f, pssIncidentClose),
	}

	cmd.Flags().String("reason", "", "Reason for incident case approval")
	cmd.MarkFlagRequired("reason")
	cmd.Flags().Int("contact", 0, "Contact ID for incident case approval")

	return cmd
}

func pssIncidentClose(service pss.PSSService, cmd *cobra.Command, args []string) error {
	closeIncidentRequest := pss.CloseIncidentCaseRequest{}
	closeIncidentRequest.Reason, _ = cmd.Flags().GetString("reason")
	closeIncidentRequest.ContactID, _ = cmd.Flags().GetInt("contact")

	var incidents []pss.IncidentCase
	for _, arg := range args {
		_, err := service.CloseIncidentCase(arg, closeIncidentRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Failed to close incident [%s]: %s", arg, err)
			continue
		}

		incident, err := service.GetIncidentCase(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving closed incident [%s]: %s", arg, err)
			continue
		}

		incidents = append(incidents, incident)
	}

	return output.CommandOutput(cmd, OutputPSSIncidentCasesProvider(incidents))
}
