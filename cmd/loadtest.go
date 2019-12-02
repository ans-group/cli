package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "loadtest",
		Short: "Commands relating to load testing service",
	}

	// Child root commands
	cmd.AddCommand(loadtestDomainRootCmd())
	cmd.AddCommand(loadtestTestRootCmd())

	return cmd
}

func outputLoadTestDomains(domains []ltaas.Domain) {
	err := Output(NewGenericOutputHandlerProvider(domains, []string{"id", "name"}))
	if err != nil {
		output.Fatalf("Failed to output domains: %s", err)
	}
}

func outputLoadTestTests(tests []ltaas.Test) {
	err := Output(NewGenericOutputHandlerProvider(tests, []string{"id", "name", "number_of_users", "duration", "protocol", "path"}))
	if err != nil {
		output.Fatalf("Failed to output tests: %s", err)
	}
}
