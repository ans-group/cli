package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssChangeRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change",
		Short: "sub-commands relating to change cases",
	}

	// Child commands
	cmd.AddCommand(pssChangeListCmd(f))
	// cmd.AddCommand(pssChangeShowCmd(f))
	// cmd.AddCommand(pssChangeCreateCmd(f))
	// cmd.AddCommand(pssChangeUpdateCmd(f))
	// cmd.AddCommand(pssChangeCloseCmd(f))

	// Child root commands
	cmd.AddCommand(pssChangeRiskRootCmd(f))
	cmd.AddCommand(pssChangeImpactRootCmd(f))

	return cmd
}

func pssChangeListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists change cases",
		Long:    "This command lists change cases",
		Example: "ans pss change list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return pssChangeList(c.PSSService(), cmd, args)
		},
	}
}

func pssChangeList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	changes, err := service.GetChangeCases(params)
	if err != nil {
		return err
	}

	return output.CommandOutput(cmd, OutputPSSChangeCasesProvider(changes))
}

// func pssChangeShowCmd(f factory.ClientFactory) *cobra.Command {
// 	return &cobra.Command{
// 		Use:     "show <change: id>...",
// 		Short:   "Shows a change",
// 		Long:    "This command shows one or more changes",
// 		Example: "ans pss change show 123",
// 		Args: func(cmd *cobra.Command, args []string) error {
// 			if len(args) < 1 {
// 				return errors.New("Missing change")
// 			}

// 			return nil
// 		},
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			c, err := f.NewClient()
// 			if err != nil {
// 				return err
// 			}

// 			return pssChangeShow(c.PSSService(), cmd, args)
// 		},
// 	}
// }

// func pssChangeShow(service pss.PSSService, cmd *cobra.Command, args []string) error {
// 	var changes []pss.Change
// 	for _, arg := range args {
// 		changeID, err := strconv.Atoi(arg)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Invalid change ID [%s]", arg)
// 			continue
// 		}

// 		change, err := service.GetChange(changeID)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Error retrieving change [%s]: %s", arg, err)
// 			continue
// 		}

// 		changes = append(changes, change)
// 	}

// 	return output.CommandOutput(cmd, OutputPSSChangesProvider(changes))
// }

// func pssChangeCreateCmd(f factory.ClientFactory) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "create",
// 		Short:   "Creates a change",
// 		Long:    "This command creates a new change",
// 		Example: "ans pss change create --subject 'example ticket' --details 'example' --author 123",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			c, err := f.NewClient()
// 			if err != nil {
// 				return err
// 			}

// 			return pssChangeCreate(c.PSSService(), cmd, args)
// 		},
// 	}

// 	// Setup flags
// 	cmd.Flags().String("subject", "", "Specifies subject for change")
// 	cmd.MarkFlagRequired("subject")
// 	cmd.Flags().String("details", "", "Specifies details for change")
// 	cmd.Flags().Int("author", 0, "Specifies author ID for change")
// 	cmd.MarkFlagRequired("author")
// 	cmd.Flags().String("priority", "Normal", "Specifies priority for change")
// 	cmd.Flags().Bool("secure", false, "Specifies whether change is secure")
// 	cmd.Flags().StringSlice("cc", []string{}, "Specifies CC email addresses for change")
// 	cmd.Flags().Bool("change-sms", false, "Specifies whether SMS updates are required")
// 	cmd.Flags().String("customer-reference", "", "Specifies customer reference for change")
// 	cmd.Flags().Int("product-id", 0, "Specifies product ID for change")
// 	cmd.Flags().String("product-name", "", "Specifies product name for change")
// 	cmd.Flags().String("product-type", "", "Specifies product type for change")

// 	return cmd
// }

// func pssChangeCreate(service pss.PSSService, cmd *cobra.Command, args []string) error {
// 	createChange := pss.CreateChangeChange{}

// 	priority, _ := cmd.Flags().GetString("priority")
// 	parsedPriority, err := pss.ChangePriorityEnum.Parse(priority)
// 	if err != nil {
// 		return err
// 	}
// 	createChange.Priority = parsedPriority

// 	if cmd.Flags().Changed("product-id") || cmd.Flags().Changed("product-name") || cmd.Flags().Changed("product-type") {
// 		createChange.Product = &pss.Product{}
// 		createChange.Product.ID, _ = cmd.Flags().GetInt("product-id")
// 		createChange.Product.Name, _ = cmd.Flags().GetString("product-name")
// 		createChange.Product.Type, _ = cmd.Flags().GetString("product-type")
// 	}

// 	createChange.Subject, _ = cmd.Flags().GetString("subject")
// 	createChange.Author.ID, _ = cmd.Flags().GetInt("author")
// 	createChange.Secure, _ = cmd.Flags().GetBool("secure")
// 	createChange.CC, _ = cmd.Flags().GetStringSlice("cc")
// 	createChange.ChangeSMS, _ = cmd.Flags().GetBool("change-sms")
// 	createChange.CustomerReference, _ = cmd.Flags().GetString("customer-reference")

// 	if cmd.Flags().Changed("details") {
// 		createChange.Details, _ = cmd.Flags().GetString("details")
// 	} else {
// 		createChange.Details, err = input.ReadInput("details")
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	changeID, err := service.CreateChange(createChange)
// 	if err != nil {
// 		return fmt.Errorf("Error creating change: %s", err)
// 	}

// 	change, err := service.GetChange(changeID)
// 	if err != nil {
// 		return fmt.Errorf("Error retrieving new change: %s", err)
// 	}

// 	return output.CommandOutput(cmd, OutputPSSChangesProvider([]pss.Change{change}))
// }

// func pssChangeUpdateCmd(f factory.ClientFactory) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:     "update <change: id>...",
// 		Short:   "Updates changes",
// 		Long:    "This command updates one or more changes",
// 		Example: "ans pss change update 123 --priority high",
// 		Args: func(cmd *cobra.Command, args []string) error {
// 			if len(args) < 1 {
// 				return errors.New("Missing change")
// 			}

// 			return nil
// 		},
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			c, err := f.NewClient()
// 			if err != nil {
// 				return err
// 			}

// 			return pssChangeUpdate(c.PSSService(), cmd, args)
// 		},
// 	}

// 	// Setup flags
// 	cmd.Flags().String("priority", "", "Specifies priority for change")
// 	cmd.Flags().String("status", "", "Specifies status for change")
// 	cmd.Flags().Bool("secure", false, "Specifies whether change is secure")
// 	cmd.Flags().Bool("read", false, "Specifies whether change is marked as read")
// 	cmd.Flags().Bool("change-sms", false, "Specifies whether SMS updates are required")
// 	cmd.Flags().Bool("archived", false, "Specifies whether change is archived")

// 	return cmd
// }

// func pssChangeUpdate(service pss.PSSService, cmd *cobra.Command, args []string) error {
// 	patchChange := pss.PatchChangeChange{}

// 	if cmd.Flags().Changed("priority") {
// 		priority, _ := cmd.Flags().GetString("priority")
// 		parsedPriority, err := pss.ChangePriorityEnum.Parse(priority)
// 		if err != nil {
// 			return err
// 		}
// 		patchChange.Priority = parsedPriority
// 	}

// 	if cmd.Flags().Changed("status") {
// 		status, _ := cmd.Flags().GetString("status")
// 		parsedStatus, err := pss.ChangeStatusEnum.Parse(status)
// 		if err != nil {
// 			return err
// 		}
// 		patchChange.Status = parsedStatus
// 	}

// 	if cmd.Flags().Changed("secure") {
// 		secure, _ := cmd.Flags().GetBool("secure")
// 		patchChange.Secure = &secure
// 	}
// 	if cmd.Flags().Changed("read") {
// 		read, _ := cmd.Flags().GetBool("read")
// 		patchChange.Read = &read
// 	}
// 	if cmd.Flags().Changed("change-sms") {
// 		changeSMS, _ := cmd.Flags().GetBool("change-sms")
// 		patchChange.ChangeSMS = &changeSMS
// 	}
// 	if cmd.Flags().Changed("archived") {
// 		archived, _ := cmd.Flags().GetBool("archived")
// 		patchChange.Archived = &archived
// 	}

// 	var changes []pss.Change

// 	for _, arg := range args {
// 		changeID, err := strconv.Atoi(arg)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Invalid change ID [%s]", arg)
// 			continue
// 		}

// 		err = service.PatchChange(changeID, patchChange)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Error updating change [%d]: %s", changeID, err)
// 			continue
// 		}

// 		change, err := service.GetChange(changeID)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Error retrieving updated change [%d]: %s", changeID, err)
// 			continue
// 		}

// 		changes = append(changes, change)
// 	}

// 	return output.CommandOutput(cmd, OutputPSSChangesProvider(changes))
// }

// func pssChangeCloseCmd(f factory.ClientFactory) *cobra.Command {
// 	return &cobra.Command{
// 		Use:     "close <change: id>...",
// 		Short:   "Closes changes",
// 		Long:    "This command closes one or more changes",
// 		Example: "ans pss change close 123",
// 		Args: func(cmd *cobra.Command, args []string) error {
// 			if len(args) < 1 {
// 				return errors.New("Missing change")
// 			}

// 			return nil
// 		},
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			c, err := f.NewClient()
// 			if err != nil {
// 				return err
// 			}

// 			return pssChangeClose(c.PSSService(), cmd, args)
// 		},
// 	}
// }

// func pssChangeClose(service pss.PSSService, cmd *cobra.Command, args []string) error {
// 	patchChange := pss.PatchChangeChange{
// 		Status: pss.ChangeStatusCompleted,
// 	}

// 	var changes []pss.Change

// 	for _, arg := range args {
// 		changeID, err := strconv.Atoi(arg)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Invalid change ID [%s]", arg)
// 			continue
// 		}

// 		err = service.PatchChange(changeID, patchChange)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Error closing change [%d]: %s", changeID, err)
// 			continue
// 		}

// 		change, err := service.GetChange(changeID)
// 		if err != nil {
// 			output.OutputWithErrorLevelf("Error retrieving updated change [%d]: %s", changeID, err)
// 			continue
// 		}

// 		changes = append(changes, change)
// 	}

// 	return output.CommandOutput(cmd, OutputPSSChangesProvider(changes))
// }
