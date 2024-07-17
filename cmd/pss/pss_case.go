package pss

import (
	"errors"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssCaseRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "case",
		Short: "sub-commands relating to cases",
		Long:  "These commands allow you to interact generically with cases. Specific case types have their own sub-commands (incident, change, problem).",
	}

	// Child commands
	cmd.AddCommand(
		pssCaseListCmd(f),
		pssCaseShowCmd(f),
	)

	// Child root commands
	cmd.AddCommand(
		pssCaseCategoryRootCmd(f),
		pssCaseUpdateRootCmd(f),
	)

	return cmd
}

func pssCaseListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists cases",
		Long:    "This command lists cases (paginated)",
		Example: "ans pss case list",
		RunE:    pssCobraRunEFunc(f, pssCaseList),
	}
}

func pssCaseList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	paginatedCases, err := service.GetCasesPaginated(params)
	if err != nil {
		return err
	}

	return output.CommandOutputPaginated(cmd, OutputPSSCasesProvider(paginatedCases.Items()), paginatedCases)
}

func pssCaseShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <case: id>...",
		Short:   "Shows a case",
		Long:    "This command shows one or more cases",
		Example: "ans pss case show CHG123456",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing case")
			}

			return nil
		},
		RunE: pssCobraRunEFunc(f, pssCaseShow),
	}
}

func pssCaseShow(service pss.PSSService, cmd *cobra.Command, args []string) error {
	var cases []pss.Case
	for _, arg := range args {
		c, err := service.GetCase(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving case [%s]: %s", arg, err)
			continue
		}

		cases = append(cases, c)
	}

	return output.CommandOutput(cmd, OutputPSSCasesProvider(cases))
}
