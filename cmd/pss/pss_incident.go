package pss

import (
	"errors"

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
	// cmd.AddCommand(pssIncidentCreateCmd(f))
	// cmd.AddCommand(pssIncidentUpdateCmd(f))
	// cmd.AddCommand(pssIncidentCloseCmd(f))

	// Child root commands
	cmd.AddCommand(pssIncidentTypeRootCmd(f))
	cmd.AddCommand(pssIncidentImpactRootCmd(f))

	return cmd
}

func pssIncidentListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists incident cases",
		Long:    "This command lists incident cases",
		Example: "ans pss incident list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return pssIncidentList(c.PSSService(), cmd, args)
		},
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return pssIncidentShow(c.PSSService(), cmd, args)
		},
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

// func pssIncidentCreateCmd(f factory.ClientFactory) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "create",
// 		Short:   "Creates a incident",
// 		Long:    "This command creates a new incident",
// 		Example: "ans pss incident create --subject 'example ticket' --details 'example' --author 123",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			c, err := f.NewClient()
// 			if err != nil {
// 				return err
// 			}

// 			return pssIncidentCreate(c.PSSService(), cmd, args)
// 		},
// 	}

// 	// Setup flags
// 	cmd.Flags().String("subject", "", "Specifies subject for incident")
// 	cmd.MarkFlagRequired("subject")
// 	cmd.Flags().String("details", "", "Specifies details for incident")
// 	cmd.Flags().Int("author", 0, "Specifies author ID for incident")
// 	cmd.MarkFlagRequired("author")
// 	cmd.Flags().String("priority", "Normal", "Specifies priority for incident")
// 	cmd.Flags().Bool("secure", false, "Specifies whether incident is secure")
// 	cmd.Flags().StringSlice("cc", []string{}, "Specifies CC email addresses for incident")
// 	cmd.Flags().Bool("incident-sms", false, "Specifies whether SMS updates are required")
// 	cmd.Flags().String("customer-reference", "", "Specifies customer reference for incident")
// 	cmd.Flags().Int("product-id", 0, "Specifies product ID for incident")
// 	cmd.Flags().String("product-name", "", "Specifies product name for incident")
// 	cmd.Flags().String("product-type", "", "Specifies product type for incident")

// 	return cmd
// }

// func pssIncidentCreate(service pss.PSSService, cmd *cobra.Command, args []string) error {
// 	createIncident := pss.CreateIncidentIncident{}

// 	priority, _ := cmd.Flags().GetString("priority")
// 	parsedPriority, err := pss.IncidentPriorityEnum.Parse(priority)
// 	if err != nil {
// 		return err
// 	}
// 	createIncident.Priority = parsedPriority

// 	if cmd.Flags().Changed("product-id") || cmd.Flags().Changed("product-name") || cmd.Flags().Changed("product-type") {
// 		createIncident.Product = &pss.Product{}
// 		createIncident.Product.ID, _ = cmd.Flags().GetInt("product-id")
// 		createIncident.Product.Name, _ = cmd.Flags().GetString("product-name")
// 		createIncident.Product.Type, _ = cmd.Flags().GetString("product-type")
// 	}

// 	createIncident.Subject, _ = cmd.Flags().GetString("subject")
// 	createIncident.Author.ID, _ = cmd.Flags().GetInt("author")
// 	createIncident.Secure, _ = cmd.Flags().GetBool("secure")
// 	createIncident.CC, _ = cmd.Flags().GetStringSlice("cc")
// 	createIncident.IncidentSMS, _ = cmd.Flags().GetBool("incident-sms")
// 	createIncident.CustomerReference, _ = cmd.Flags().GetString("customer-reference")

// 	if cmd.Flags().Changed("details") {
// 		createIncident.Details, _ = cmd.Flags().GetString("details")
// 	} else {
// 		createIncident.Details, err = input.ReadInput("details")
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	incidentID, err := service.CreateIncident(createIncident)
// 	if err != nil {
// 		return fmt.Errorf("Error creating incident: %s", err)
// 	}

// 	incident, err := service.GetIncident(incidentID)
// 	if err != nil {
// 		return fmt.Errorf("Error retrieving new incident: %s", err)
// 	}

// 	return output.CommandOutput(cmd, OutputPSSIncidentsProvider([]pss.Incident{incident}))
// }

// func pssIncidentUpdateCmd(f factory.ClientFactory) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "update <incident: id>...",
// 		Short:   "Updates incidents",
// 		Long:    "This command updates one or more incidents",
// 		Example: "ans pss incident update 123 --priority high",
// 		Args: func(cmd *cobra.Command, args []string) error {
// 			if len(args) < 1 {
// 				return errors.New("Missing incident")
// 			}

// 			return nil
// 		},
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			c, err := f.NewClient()
// 			if err != nil {
// 				return err
// 			}

// 			return pssIncidentUpdate(c.PSSService(), cmd, args)
// 		},
// 	}

// 	// Setup flags
// 	cmd.Flags().String("priority", "", "Specifies priority for incident")
// 	cmd.Flags().String("status", "", "Specifies status for incident")
// 	cmd.Flags().Bool("secure", false, "Specifies whether incident is secure")
// 	cmd.Flags().Bool("read", false, "Specifies whether incident is marked as read")
// 	cmd.Flags().Bool("incident-sms", false, "Specifies whether SMS updates are required")
// 	cmd.Flags().Bool("archived", false, "Specifies whether incident is archived")

// 	return cmd
// }

// func pssIncidentUpdate(service pss.PSSService, cmd *cobra.Command, args []string) error {
// 	patchIncident := pss.PatchIncidentIncident{}

// 	if cmd.Flags().Changed("priority") {
// 		priority, _ := cmd.Flags().GetString("priority")
// 		parsedPriority, err := pss.IncidentPriorityEnum.Parse(priority)
// 		if err != nil {
// 			return err
// 		}
// 		patchIncident.Priority = parsedPriority
// 	}

// 	if cmd.Flags().Changed("status") {
// 		status, _ := cmd.Flags().GetString("status")
// 		parsedStatus, err := pss.IncidentStatusEnum.Parse(status)
// 		if err != nil {
// 			return err
// 		}
// 		patchIncident.Status = parsedStatus
// 	}

// 	if cmd.Flags().Changed("secure") {
// 		secure, _ := cmd.Flags().GetBool("secure")
// 		patchIncident.Secure = &secure
// 	}
// 	if cmd.Flags().Changed("read") {
// 		read, _ := cmd.Flags().GetBool("read")
// 		patchIncident.Read = &read
// 	}
// 	if cmd.Flags().Changed("incident-sms") {
// 		incidentSMS, _ := cmd.Flags().GetBool("incident-sms")
// 		patchIncident.IncidentSMS = &incidentSMS
// 	}
// 	if cmd.Flags().Changed("archived") {
// 		archived, _ := cmd.Flags().GetBool("archived")
// 		patchIncident.Archived = &archived
// 	}

// 	var incidents []pss.Incident

// 	for _, arg := range args {
// 		incidentID, err := strconv.Atoi(arg)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Invalid incident ID [%s]", arg)
// 			continue
// 		}

// 		err = service.PatchIncident(incidentID, patchIncident)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Error updating incident [%d]: %s", incidentID, err)
// 			continue
// 		}

// 		incident, err := service.GetIncident(incidentID)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Error retrieving updated incident [%d]: %s", incidentID, err)
// 			continue
// 		}

// 		incidents = append(incidents, incident)
// 	}

// 	return output.CommandOutput(cmd, OutputPSSIncidentsProvider(incidents))
// }

// func pssIncidentCloseCmd(f factory.ClientFactory) *cobra.Command {
// 	return &cobra.Command{
// 		Use:     "close <incident: id>...",
// 		Short:   "Closes incidents",
// 		Long:    "This command closes one or more incidents",
// 		Example: "ans pss incident close 123",
// 		Args: func(cmd *cobra.Command, args []string) error {
// 			if len(args) < 1 {
// 				return errors.New("Missing incident")
// 			}

// 			return nil
// 		},
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			c, err := f.NewClient()
// 			if err != nil {
// 				return err
// 			}

// 			return pssIncidentClose(c.PSSService(), cmd, args)
// 		},
// 	}
// }

// func pssIncidentClose(service pss.PSSService, cmd *cobra.Command, args []string) error {
// 	patchIncident := pss.PatchIncidentIncident{
// 		Status: pss.IncidentStatusCompleted,
// 	}

// 	var incidents []pss.Incident

// 	for _, arg := range args {
// 		incidentID, err := strconv.Atoi(arg)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Invalid incident ID [%s]", arg)
// 			continue
// 		}

// 		err = service.PatchIncident(incidentID, patchIncident)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Error closing incident [%d]: %s", incidentID, err)
// 			continue
// 		}

// 		incident, err := service.GetIncident(incidentID)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Error retrieving updated incident [%d]: %s", incidentID, err)
// 			continue
// 		}

// 		incidents = append(incidents, incident)
// 	}

// 	return output.CommandOutput(cmd, OutputPSSIncidentsProvider(incidents))
// }
