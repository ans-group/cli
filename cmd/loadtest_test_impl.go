package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
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
	cmd.AddCommand(loadtestTestDeleteCmd())

	// Child root commands
	cmd.AddCommand(loadtestTestJobRootCmd())

	return cmd
}

func loadtestTestListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists tests",
		Long:    "This command lists tests",
		Example: "ukfast loadtest test list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return loadtestTestList(getClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestTestList(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	tests, err := service.GetTests(params)
	if err != nil {
		return fmt.Errorf("Error retrieving tests: %s", err)
	}

	return outputLoadTestTests(tests)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return loadtestTestShow(getClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestTestShow(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	var tests []ltaas.Test
	for _, arg := range args {
		test, err := service.GetTest(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving test [%s]: %s", arg, err)
			continue
		}

		tests = append(tests, test)
	}

	return outputLoadTestTests(tests)
}

func loadtestTestDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <test: id>...",
		Short:   "Deletes a test",
		Long:    "This command deletes one or more tests",
		Example: "ukfast loadtest test delete 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing test")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return loadtestTestDelete(getClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestTestDelete(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	var tests []ltaas.Test
	for _, arg := range args {
		test, err := service.GetTest(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing test [%s]: %s", arg, err)
			continue
		}

		tests = append(tests, test)
	}

	return outputLoadTestTests(tests)
}
