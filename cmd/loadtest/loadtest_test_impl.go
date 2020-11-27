package loadtest

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	flaghelper "github.com/ukfast/cli/internal/pkg/helper/flag"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestTestRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "sub-commands relating to tests",
	}

	// Child commands
	cmd.AddCommand(loadtestTestListCmd(f))
	cmd.AddCommand(loadtestTestShowCmd(f))
	cmd.AddCommand(loadtestTestCreateCmd(f))
	cmd.AddCommand(loadtestTestDeleteCmd(f))

	return cmd
}

func loadtestTestListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists tests",
		Long:    "This command lists tests",
		Example: "ukfast loadtest test list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadtestTestList(c.LTaaSService(), cmd, args)
		},
	}
}

func loadtestTestList(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	params, err := flaghelper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	tests, err := service.GetTests(params)
	if err != nil {
		return fmt.Errorf("Error retrieving tests: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadTestTestsProvider(tests))
}

func loadtestTestShowCmd(f factory.ClientFactory) *cobra.Command {
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
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadtestTestShow(c.LTaaSService(), cmd, args)
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

	return output.CommandOutput(cmd, OutputLoadTestTestsProvider(tests))
}

func loadtestTestCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a test",
		Long:    "This command creates a test ",
		Example: "ukfast loadtest test create",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadtestTestCreate(c.LTaaSService(), cmd, args)
		},
	}

	cmd.Flags().String("domain-id", "", "ID for domain")
	cmd.MarkFlagRequired("domain-id")
	cmd.Flags().String("name", "", "Name for test")
	cmd.MarkFlagRequired("name")
	cmd.Flags().String("scenario-id", "", "ID for scenario to use for test")
	cmd.MarkFlagRequired("scenario-id")
	cmd.Flags().String("script-id", "", "ID for script")
	cmd.Flags().String("protocol", "", "Protocol to use")
	cmd.MarkFlagRequired("protocol")
	cmd.Flags().String("path", "", "URI path to use, e.g. /blog")
	cmd.Flags().Int("number-of-users", 0, "Number of users for test")
	cmd.MarkFlagRequired("number-of-users")
	cmd.Flags().String("duration", "", "Duration of test, e.g. 5m30s")
	cmd.MarkFlagRequired("duration")
	cmd.Flags().String("authorisation-agreement-version", "", "Version of authorisation agreement, e.g. v1.0")
	cmd.MarkFlagRequired("authorisation-agreement-version")
	cmd.Flags().String("authorisation-name", "", "Name of person who authorised the test")
	cmd.MarkFlagRequired("authorisation-name")
	cmd.Flags().String("authorisation-position", "", "Position of person who authorised the test")
	cmd.MarkFlagRequired("authorisation-position")
	cmd.Flags().String("authorisation-company", "", "Company of person who authorised the test")
	cmd.MarkFlagRequired("authorisation-company")

	return cmd
}

func loadtestTestCreate(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	createRequest := ltaas.CreateTestRequest{}
	createRequest.DomainID, _ = cmd.Flags().GetString("domain-id")
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.ScenarioID, _ = cmd.Flags().GetString("scenario-id")
	createRequest.ScriptID, _ = cmd.Flags().GetString("script-id")
	createRequest.Path, _ = cmd.Flags().GetString("path")
	createRequest.NumberOfUsers, _ = cmd.Flags().GetInt("number-of-users")

	protocol, _ := cmd.Flags().GetString("protocol")
	parsedProtocol, err := ltaas.ParseTestProtocol(protocol)
	if err != nil {
		return nil
	}
	createRequest.Protocol = parsedProtocol

	duration, _ := cmd.Flags().GetString("duration")
	parsedDuration, err := ltaas.ParseTestDuration(duration)
	if err != nil {
		return nil
	}
	createRequest.Duration = parsedDuration

	authorization := ltaas.CreateTestAuthorisation{}
	authorization.AgreementVersion, _ = cmd.Flags().GetString("authorisation-agreement-version")
	authorization.Name, _ = cmd.Flags().GetString("authorisation-name")
	authorization.Position, _ = cmd.Flags().GetString("authorisation-position")
	authorization.Company, _ = cmd.Flags().GetString("authorisation-company")
	createRequest.Authorisation = authorization

	testID, err := service.CreateTest(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating test: %s", err)
	}

	test, err := service.GetTest(testID)
	if err != nil {
		return fmt.Errorf("Error retrieving new test [%s]: %s", testID, err)
	}

	return output.CommandOutput(cmd, OutputLoadTestTestsProvider([]ltaas.Test{test}))
}

func loadtestTestDeleteCmd(f factory.ClientFactory) *cobra.Command {
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
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			loadtestTestDelete(c.LTaaSService(), cmd, args)
			return nil
		},
	}
}

func loadtestTestDelete(service ltaas.LTaaSService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		err := service.DeleteTest(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing test [%s]: %s", arg, err)
			continue
		}
	}
}
