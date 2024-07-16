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

func pssCaseUpdateRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "sub-commands relating to case updates",
	}

	// Child root commands
	cmd.AddCommand(pssCaseUpdateListCmd(f))
	cmd.AddCommand(pssCaseUpdateShowCmd(f))

	return cmd
}

func pssCaseUpdateListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <case: id>",
		Short:   "List updates for a case",
		Long:    "This command lists updates for a case",
		Example: "ans pss case update list CHG123456",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing case")
			}

			return nil
		},
		RunE: pssCobraRunEFunc(f, pssCaseUpdateList),
	}
}

func pssCaseUpdateList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	updates, err := service.GetCaseUpdates(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving case updates: %s", err)
	}

	return output.CommandOutput(cmd, OutputPSSCaseUpdatesProvider(updates))
}

func pssCaseUpdateShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <case: id>...",
		Short:   "Shows a case update",
		Long:    "This command shows one or more case updates",
		Example: "ans pss case update show CHG123456 9ddf3546-1d14-4604-acfa-aebcb6a32ec9",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing case")
			}
			if len(args) < 2 {
				return errors.New("Missing case update")
			}

			return nil
		},
		RunE: pssCobraRunEFunc(f, pssCaseUpdateShow),
	}
}

func pssCaseUpdateShow(service pss.PSSService, cmd *cobra.Command, args []string) error {
	var cases []pss.CaseUpdate
	for _, arg := range args[1:] {
		c, err := service.GetCaseUpdate(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving case update [%s]: %s", arg, err)
			continue
		}

		cases = append(cases, c)
	}

	return output.CommandOutput(cmd, OutputPSSCaseUpdatesProvider(cases))
}
