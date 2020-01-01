package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestTestRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "sub-commands relating to tests",
	}

	// Child commands
	cmd.AddCommand(loadtestTestListCmd())
	cmd.AddCommand(loadtestTestShowCmd())

	return cmd
}

func loadtestTestListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists tests",
		Long:    "This command lists tests",
		Example: "ukfast loadtest test list",
		Run: func(cmd *cobra.Command, args []string) {
			loadtestTestList(getClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestTestList(service ltaas.LTaaSService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	tests, err := service.GetTests(params)
	if err != nil {
		output.Fatalf("Error retrieving tests: %s", err)
		return
	}

	outputLoadTestTests(tests)
}

func loadtestTestShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <test: id>...",
		Short:   "Shows a test",
		Long:    "This command shows one or more tests",
		Example: "ukfast loadtest test show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing test")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			loadtestTestShow(getClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestTestShow(service ltaas.LTaaSService, cmd *cobra.Command, args []string) {
	var tests []ltaas.Test
	for _, arg := range args {
		test, err := service.GetTest(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving test [%s]: %s", arg, err)
			continue
		}

		tests = append(tests, test)
	}

	outputLoadTestTests(tests)
}
