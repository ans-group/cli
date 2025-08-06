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

func pssChangeRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change",
		Short: "sub-commands relating to change cases",
	}

	// Child commands
	cmd.AddCommand(pssChangeListCmd(f))
	cmd.AddCommand(pssChangeShowCmd(f))
	cmd.AddCommand(pssChangeCreateCmd(f))
	cmd.AddCommand(pssChangeApproveCmd(f))

	// Child root commands
	cmd.AddCommand(pssChangeRiskRootCmd(f))
	cmd.AddCommand(pssChangeImpactRootCmd(f))

	// Additional root commands (generic case commands)
	cmd.AddCommand(pssCaseUpdateRootCmd(f))
	cmd.AddCommand(pssCaseCategoryRootCmd(f))

	return cmd
}

func pssChangeListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists change cases",
		Long:    "This command lists change cases (paginated)",
		Example: "ans pss change list",
		RunE:    pssCobraRunEFunc(f, pssChangeList),
	}
}

func pssChangeList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	paginatedChanges, err := service.GetChangeCasesPaginated(params)
	if err != nil {
		return err
	}

	return output.CommandOutputPaginated(cmd, ChangeCaseCollection(paginatedChanges.Items()), paginatedChanges)
}

func pssChangeShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <change: id>...",
		Short:   "Shows an change",
		Long:    "This command shows one or more changes",
		Example: "ans pss change show CHG123456",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing change")
			}

			return nil
		},
		RunE: pssCobraRunEFunc(f, pssChangeShow),
	}
}

func pssChangeShow(service pss.PSSService, cmd *cobra.Command, args []string) error {
	var changes []pss.ChangeCase
	for _, arg := range args {
		change, err := service.GetChangeCase(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving change [%s]: %s", arg, err)
			continue
		}

		changes = append(changes, change)
	}

	return output.CommandOutput(cmd, ChangeCaseCollection(changes))
}

func pssChangeCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a change",
		Long:    "This command creates a new change",
		Example: "ans pss change create --title 'test change' --description 'test change' --risk Low --category 70b67a49-7ace-4146-a295-11590a0b1203 --supported-service 1b684067-a587-4997-acb9-da2f4cb7be81 --impact Low --reason 'test reason'",
		RunE:    pssCobraRunEFunc(f, pssChangeCreate),
	}

	// Setup flags
	cmd.Flags().String("title", "", "Specifies the title for change case")
	_ = cmd.MarkFlagRequired("title")
	cmd.Flags().String("description", "", "Specifies the description for change case")
	_ = cmd.MarkFlagRequired("description")
	cmd.Flags().String("risk", "", "Specifies the risk of change case")
	_ = cmd.MarkFlagRequired("risk")
	cmd.Flags().String("category", "", "Category ID for change case")
	_ = cmd.MarkFlagRequired("category")
	cmd.Flags().String("supported-service", "", "Supported service ID for change case")
	_ = cmd.MarkFlagRequired("supported-service")
	cmd.Flags().String("impact", "", "Impact for change case")
	_ = cmd.MarkFlagRequired("impact")
	cmd.Flags().String("reason", "", "Reason for change case")
	_ = cmd.MarkFlagRequired("reason")
	cmd.Flags().Bool("security", false, "Specifies whether change case is a security change")
	cmd.Flags().String("customer-reference", "", "Specifies the customer reference for change case")

	return cmd
}

func pssChangeCreate(service pss.PSSService, cmd *cobra.Command, args []string) error {
	createChangeCase := pss.CreateChangeCaseRequest{}

	createChangeCase.Title, _ = cmd.Flags().GetString("title")
	createChangeCase.Description, _ = cmd.Flags().GetString("description")
	createChangeCase.IsSecurity, _ = cmd.Flags().GetBool("security")
	createChangeCase.CustomerReference, _ = cmd.Flags().GetString("customer-reference")
	createChangeCase.Reason, _ = cmd.Flags().GetString("reason")
	createChangeCase.CategoryID, _ = cmd.Flags().GetString("category")
	createChangeCase.SupportedServiceID, _ = cmd.Flags().GetString("supported-service")

	changeCaseRisk, _ := cmd.Flags().GetString("risk")
	parsedChangeCaseRisk, err := pss.ChangeCaseRiskEnum.Parse(changeCaseRisk)
	if err != nil {
		return err
	}
	createChangeCase.Risk = parsedChangeCaseRisk

	changeCaseImpact, _ := cmd.Flags().GetString("impact")
	parsedChangeCaseImpact, err := pss.ChangeCaseImpactEnum.Parse(changeCaseImpact)
	if err != nil {
		return err
	}
	createChangeCase.Impact = parsedChangeCaseImpact

	changeID, err := service.CreateChangeCase(createChangeCase)
	if err != nil {
		return fmt.Errorf("error creating change: %s", err)
	}

	change, err := service.GetChangeCase(changeID)
	if err != nil {
		return fmt.Errorf("error retrieving new change: %s", err)
	}

	return output.CommandOutput(cmd, ChangeCaseCollection([]pss.ChangeCase{change}))
}

func pssChangeApproveCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "approve <change: id>...",
		Short:   "Approves a change",
		Long:    "This command approves one or more changes",
		Example: "ans pss change approve CHG123456",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing change")
			}

			return nil
		},
		RunE: pssCobraRunEFunc(f, pssChangeApprove),
	}

	cmd.Flags().String("reason", "", "Reason for change case approval")
	cmd.Flags().Int("contact", 0, "Contact ID for incident case approval")

	return cmd
}

func pssChangeApprove(service pss.PSSService, cmd *cobra.Command, args []string) error {
	approveChangeRequest := pss.ApproveChangeCaseRequest{}
	approveChangeRequest.Reason, _ = cmd.Flags().GetString("reason")
	approveChangeRequest.ContactID, _ = cmd.Flags().GetInt("contact")

	var changes []pss.ChangeCase
	for _, arg := range args {
		_, err := service.ApproveChangeCase(arg, approveChangeRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Failed to approve change [%s]: %s", arg, err)
			continue
		}

		change, err := service.GetChangeCase(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving approved change [%s]: %s", arg, err)
			continue
		}

		changes = append(changes, change)
	}

	return output.CommandOutput(cmd, ChangeCaseCollection(changes))
}
