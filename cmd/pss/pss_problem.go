package pss

import (
	"errors"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssProblemRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "problem",
		Short: "sub-commands relating to problem cases",
	}

	// Child commands
	cmd.AddCommand(pssProblemListCmd(f))
	cmd.AddCommand(pssProblemShowCmd(f))

	// Additional root commands (generic case commands)
	cmd.AddCommand(pssCaseUpdateRootCmd(f))
	cmd.AddCommand(pssCaseCategoryRootCmd(f))

	return cmd
}

func pssProblemListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists problem cases",
		Long:    "This command lists problem cases (paginated)",
		Example: "ans pss problem list",
		RunE:    pssCobraRunEFunc(f, pssProblemList),
	}
}

func pssProblemList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	paginatedProblems, err := service.GetProblemCasesPaginated(params)
	if err != nil {
		return err
	}

	return output.CommandOutputPaginated(cmd, OutputPSSProblemCasesProvider(paginatedProblems.Items()), paginatedProblems)
}

func pssProblemShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <problem: id>...",
		Short:   "Shows an problem",
		Long:    "This command shows one or more problems",
		Example: "ans pss problem show PRB123456",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing problem")
			}

			return nil
		},
		RunE: pssCobraRunEFunc(f, pssProblemShow),
	}
}

func pssProblemShow(service pss.PSSService, cmd *cobra.Command, args []string) error {
	var problems []pss.ProblemCase
	for _, arg := range args {
		problem, err := service.GetProblemCase(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving problem [%s]: %s", arg, err)
			continue
		}

		problems = append(problems, problem)
	}

	return output.CommandOutput(cmd, OutputPSSProblemCasesProvider(problems))
}
