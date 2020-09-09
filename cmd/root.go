package cmd

import (
	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	accountcmd "github.com/ukfast/cli/cmd/account"
	billingcmd "github.com/ukfast/cli/cmd/billing"
	ddosxcmd "github.com/ukfast/cli/cmd/ddosx"
	draascmd "github.com/ukfast/cli/cmd/draas"
	ecloudcmd "github.com/ukfast/cli/cmd/ecloud"
	ecloudflexcmd "github.com/ukfast/cli/cmd/ecloudflex"
	loadtestcmd "github.com/ukfast/cli/cmd/loadtest"
	psscmd "github.com/ukfast/cli/cmd/pss"
	registrarcmd "github.com/ukfast/cli/cmd/registrar"
	safednscmd "github.com/ukfast/cli/cmd/safedns"
	sharedexchangecmd "github.com/ukfast/cli/cmd/sharedexchange"
	sslcmd "github.com/ukfast/cli/cmd/ssl"
	storagecmd "github.com/ukfast/cli/cmd/storage"
	threatmonitoringcmd "github.com/ukfast/cli/cmd/threatmonitoring"
	"github.com/ukfast/cli/internal/pkg/build"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
)

var appVersion string
var defaultConfigFile string

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
	rootCmd.PersistentFlags().StringP("output", "o", "", "output type {table, json, jsonpath, template, value, csv, list}, with optional argument provided as 'outputname=outputargument'")
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
	rootCmd.AddCommand(ConfigRootCmd(fs))
	rootCmd.AddCommand(CompletionRootCmd())
	rootCmd.AddCommand(accountcmd.AccountRootCmd(clientFactory))
	rootCmd.AddCommand(billingcmd.BillingRootCmd(clientFactory))
	rootCmd.AddCommand(ddosxcmd.DDoSXRootCmd(clientFactory, fs))
	rootCmd.AddCommand(draascmd.DRaaSRootCmd(clientFactory))
	rootCmd.AddCommand(ecloudcmd.ECloudRootCmd(clientFactory))
	rootCmd.AddCommand(ecloudflexcmd.ECloudFlexRootCmd(clientFactory))
	rootCmd.AddCommand(loadtestcmd.LoadTestRootCmd(clientFactory))
	rootCmd.AddCommand(psscmd.PSSRootCmd(clientFactory, fs))
	rootCmd.AddCommand(registrarcmd.RegistrarRootCmd(clientFactory))
	rootCmd.AddCommand(safednscmd.SafeDNSRootCmd(clientFactory))
	rootCmd.AddCommand(sharedexchangecmd.SharedExchangeRootCmd(clientFactory))
	rootCmd.AddCommand(sslcmd.SSLRootCmd(clientFactory, fs))
	rootCmd.AddCommand(storagecmd.StorageRootCmd(clientFactory))
	rootCmd.AddCommand(threatmonitoringcmd.ThreatMonitoringRootCmd(clientFactory))

	if err := rootCmd.Execute(); err != nil {
		output.Fatal(err.Error())
	}

	output.ExitWithErrorLevel()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("ukf")
	viper.AutomaticEnv() // read in environment variables that match

	var configFilePath string
	configFile := rootCmd.Flags().Changed("config")
	if configFile {
		configFilePath, _ = rootCmd.Flags().GetString("config")
		// Use config file from the flag.
		viper.SetConfigFile(configFilePath)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			output.Fatal(err.Error())
		}

		// Search config in home directory with name ".ukfast" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ukfast")
		defaultConfigFile = fmt.Sprintf("%s/.ukfast.yml", home)
	}

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if configFile && err != nil {
		output.Fatalf("Failed to read config from file '%s': %s", configFilePath, err.Error())
	}
}
