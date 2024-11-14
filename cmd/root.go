package cmd

import (
	accountcmd "github.com/ans-group/cli/cmd/account"
	billingcmd "github.com/ans-group/cli/cmd/billing"
	cloudflarecmd "github.com/ans-group/cli/cmd/cloudflare"
	configcmd "github.com/ans-group/cli/cmd/config"
	ddosxcmd "github.com/ans-group/cli/cmd/ddosx"
	draascmd "github.com/ans-group/cli/cmd/draas"
	ecloudcmd "github.com/ans-group/cli/cmd/ecloud"
	loadbalancercmd "github.com/ans-group/cli/cmd/loadbalancer"
	psscmd "github.com/ans-group/cli/cmd/pss"
	registrarcmd "github.com/ans-group/cli/cmd/registrar"
	safednscmd "github.com/ans-group/cli/cmd/safedns"
	sslcmd "github.com/ans-group/cli/cmd/ssl"
	storagecmd "github.com/ans-group/cli/cmd/storage"
	"github.com/ans-group/cli/internal/pkg/build"
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/config"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/logging"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var appVersion string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ans",
	Short:   "Utility for manipulating ANS services",
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
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.ans.yml)")
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
	connectionFactory := connection.NewDefaultConnectionFactory(
		connection.WithDefaultConnectionUserAgent("ans-cli"),
	)
	clientFactory := factory.NewANSClientFactory(connectionFactory)

	// Child commands
	rootCmd.AddCommand(updateCmd())

	// Child root commands
	rootCmd.AddCommand(configcmd.ConfigRootCmd(fs))
	rootCmd.AddCommand(CompletionRootCmd())
	rootCmd.AddCommand(rawCmd(connectionFactory))
	rootCmd.AddCommand(accountcmd.AccountRootCmd(clientFactory))
	rootCmd.AddCommand(billingcmd.BillingRootCmd(clientFactory))
	rootCmd.AddCommand(ddosxcmd.DDoSXRootCmd(clientFactory, fs))
	rootCmd.AddCommand(draascmd.DRaaSRootCmd(clientFactory))
	rootCmd.AddCommand(ecloudcmd.ECloudRootCmd(clientFactory, fs))
	rootCmd.AddCommand(loadbalancercmd.LoadBalancerRootCmd(clientFactory, fs))
	rootCmd.AddCommand(cloudflarecmd.CloudflareRootCmd(clientFactory))
	rootCmd.AddCommand(psscmd.PSSRootCmd(clientFactory, fs))
	rootCmd.AddCommand(registrarcmd.RegistrarRootCmd(clientFactory))
	rootCmd.AddCommand(safednscmd.SafeDNSRootCmd(clientFactory))
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

	if config.GetBool("api_debug") {
		logging.SetLogger(&output.DebugLogger{})
	}
}
