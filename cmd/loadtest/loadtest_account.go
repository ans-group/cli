package loadtest

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestAccountRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "sub-commands relating to accounts",
	}

	// Child commands
	cmd.AddCommand(loadtestAccountCreateCmd(f))

	return cmd
}

func loadtestAccountCreateCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "create",
		Short:   "Creates a account",
		Long:    "This command creates a account ",
		Example: "ukfast loadtest account create",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadtestAccountCreate(c.LTaaSService(), cmd, args)
		},
	}
}

func loadtestAccountCreate(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	_, err := service.CreateAccount()
	if err != nil {
		return fmt.Errorf("Error creating account: %s", err)
	}

	var tests []ltaas.Test
	for _, arg := range args {
		test, err := service.GetTest(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving test [%s]: %s", arg, err)
			continue
		}

		tests = append(tests, test)
	}

	return output.CommandOutput(cmd, OutputLoadTestTestsProvider(tests))
}
