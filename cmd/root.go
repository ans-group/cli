package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	accountcmd "github.com/ukfast/cli/cmd/account"
	billingcmd "github.com/ukfast/cli/cmd/billing"
	cloudflarecmd "github.com/ukfast/cli/cmd/cloudflare"
	configcmd "github.com/ukfast/cli/cmd/config"
	ddosxcmd "github.com/ukfast/cli/cmd/ddosx"
	draascmd "github.com/ukfast/cli/cmd/draas"
	ecloudcmd "github.com/ukfast/cli/cmd/ecloud"
	ecloudflexcmd "github.com/ukfast/cli/cmd/ecloudflex"
	loadbalancercmd "github.com/ukfast/cli/cmd/loadbalancer"
	psscmd "github.com/ukfast/cli/cmd/pss"
	registrarcmd "github.com/ukfast/cli/cmd/registrar"
	safednscmd "github.com/ukfast/cli/cmd/safedns"
	sharedexchangecmd "github.com/ukfast/cli/cmd/sharedexchange"
	sslcmd "github.com/ukfast/cli/cmd/ssl"
	storagecmd "github.com/ukfast/cli/cmd/storage"
	"github.com/ukfast/cli/internal/pkg/build"
	"github.com/ukfast/cli/internal/pkg/config"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
)

var appVersion string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ukfast",
	Short:   "Utility for manipulating UKFast services",
	Version: "UNKNOWN",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(build build.BuildInfo) {
	appVersion = build.Version
	rootCmd.Version = build.String()
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	// Global flags
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.ukfast.yml)")
	rootCmd.PersistentFlags().String("context", "", "specific context to use")
	rootCmd.PersistentFlags().StringP("output", "o", "", "output type {table, json, yaml, jsonpath, template, value, csv, list}, with optional argument provided as 'outputname=outputargument'")
	rootCmd.PersistentFlags().StringP("format", "f", "", "")
	rootCmd.PersistentFlags().MarkDeprecated("format", "please use --output/-o instead")
	rootCmd.PersistentFlags().String("outputtemplate", "", "output Go template (used with 'template' format), e.g. 'Name: {{ .Name }}'")
	rootCmd.PersistentFlags().MarkDeprecated("outputtemplate", "please use --output/-o flag args instead (see documentation)")
	rootCmd.PersistentFlags().String("sort", "", "output sorting, e.g. 'name', 'name:asc', 'name:desc'")
	rootCmd.PersistentFlags().StringSlice("property", []string{}, "property to output (used with several formats), can be repeated")
	rootCmd.PersistentFlags().StringArray("filter", []string{}, "filter for list commands, can be repeated, e.g. 'property=somevalue', 'property:gt=3', 'property=valu*'")
	rootCmd.PersistentFlags().Int("page", 0, "page to retrieve for paginated requests")

	cobra.OnInitialize(initConfig)
	fs := afero.NewOsFs()
	clientFactory := factory.NewUKFastClientFactory(
		factory.WithUserAgent("ukfast-cli"),
	)

	// Child commands
	rootCmd.AddCommand(updateCmd())

	// Child root commands
	rootCmd.AddCommand(configcmd.ConfigRootCmd(fs))
	rootCmd.AddCommand(CompletionRootCmd())
	rootCmd.AddCommand(accountcmd.AccountRootCmd(clientFactory))
	rootCmd.AddCommand(billingcmd.BillingRootCmd(clientFactory))
	rootCmd.AddCommand(ddosxcmd.DDoSXRootCmd(clientFactory, fs))
	rootCmd.AddCommand(draascmd.DRaaSRootCmd(clientFactory))
	rootCmd.AddCommand(ecloudcmd.ECloudRootCmd(clientFactory, fs))
	rootCmd.AddCommand(ecloudflexcmd.ECloudFlexRootCmd(clientFactory))
	rootCmd.AddCommand(loadbalancercmd.LoadBalancerRootCmd(clientFactory, fs))
	rootCmd.AddCommand(cloudflarecmd.CloudflareRootCmd(clientFactory))
	rootCmd.AddCommand(psscmd.PSSRootCmd(clientFactory, fs))
	rootCmd.AddCommand(registrarcmd.RegistrarRootCmd(clientFactory))
	rootCmd.AddCommand(safednscmd.SafeDNSRootCmd(clientFactory))
	rootCmd.AddCommand(sharedexchangecmd.SharedExchangeRootCmd(clientFactory))
	rootCmd.AddCommand(sslcmd.SSLRootCmd(clientFactory, fs))
	rootCmd.AddCommand(storagecmd.StorageRootCmd(clientFactory))

	if err := rootCmd.Execute(); err != nil {
		output.Fatal(err.Error())
	}

	output.ExitWithErrorLevel()
}

// initConfig initialises config
func initConfig() {
	configPath, _ := rootCmd.Flags().GetString("config")
	err := config.Init(configPath)
	if err != nil {
		output.Fatalf("Failed to initialise config: %s", err.Error())
	}

	if rootCmd.Flags().Changed("context") {
		contextName, _ := rootCmd.Flags().GetString("context")
		err := config.SwitchCurrentContext(contextName)
		if err != nil {
			output.Fatalf("Failed to set context: %s", err)
		}
	}
}
