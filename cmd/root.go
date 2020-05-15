package cmd

import (
	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	accountcmd "github.com/ukfast/cli/cmd/account"
	ddosxcmd "github.com/ukfast/cli/cmd/ddosx"
	ecloudcmd "github.com/ukfast/cli/cmd/ecloud"
	loadtestcmd "github.com/ukfast/cli/cmd/loadtest"
	psscmd "github.com/ukfast/cli/cmd/pss"
	registrarcmd "github.com/ukfast/cli/cmd/registrar"
	safednscmd "github.com/ukfast/cli/cmd/safedns"
	sslcmd "github.com/ukfast/cli/cmd/ssl"
	storagecmd "github.com/ukfast/cli/cmd/storage"
	"github.com/ukfast/cli/internal/pkg/build"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
)

var flagConfig string
var flagFormat string
var flagOutputTemplate string
var flagJSONPath string
var flagSort string
var flagProperty []string
var flagFilter []string
var fs afero.Fs
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
	rootCmd.PersistentFlags().StringVar(&flagConfig, "config", "", "config file (default is $HOME/.ukfast.yml)")
	rootCmd.PersistentFlags().StringVarP(&flagFormat, "format", "f", "", "output format {table, json, jsonpath, template, value, csv, list}")
	rootCmd.PersistentFlags().StringVar(&flagOutputTemplate, "gotemplate", "", "output Go template (used with 'template' format), e.g. 'Name: {{ .Name }}'")
	rootCmd.PersistentFlags().StringVar(&flagJSONPath, "jsonpath", "", "JSON path query (used with 'jsonpath' format)")
	rootCmd.PersistentFlags().StringVar(&flagSort, "sort", "", "output sorting, e.g. 'name', 'name:asc', 'name:desc'")
	rootCmd.PersistentFlags().StringSliceVar(&flagProperty, "property", []string{}, "property to output (used with several formats), can be repeated")
	rootCmd.PersistentFlags().StringArrayVar(&flagFilter, "filter", []string{}, "filter for list commands, can be repeated, e.g. 'property=somevalue', 'property:gt=3', 'property=valu*'")

	cobra.OnInitialize(initConfig)
	fs = afero.NewOsFs()
	clientFactory := factory.NewUKFastClientFactory(
		factory.WithUserAgent("ukfast-cli"),
	)

	// Child commands
	rootCmd.AddCommand(updateCmd())

	// Child root commands
	rootCmd.AddCommand(ConfigRootCmd(fs))
	rootCmd.AddCommand(CompletionRootCmd())
	rootCmd.AddCommand(accountcmd.AccountRootCmd(clientFactory))
	rootCmd.AddCommand(ddosxcmd.DDoSXRootCmd(clientFactory, fs))
	rootCmd.AddCommand(ecloudcmd.ECloudRootCmd(clientFactory))
	rootCmd.AddCommand(loadtestcmd.LoadTestRootCmd(clientFactory))
	rootCmd.AddCommand(psscmd.PSSRootCmd(clientFactory, fs))
	rootCmd.AddCommand(registrarcmd.RegistrarRootCmd(clientFactory))
	rootCmd.AddCommand(safednscmd.SafeDNSRootCmd(clientFactory))
	rootCmd.AddCommand(sslcmd.SSLRootCmd(clientFactory))
	rootCmd.AddCommand(storagecmd.StorageRootCmd(clientFactory))

	if err := rootCmd.Execute(); err != nil {
		output.Fatal(err.Error())
	}

	output.ExitWithErrorLevel()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if flagConfig != "" {
		// Use config file from the flag.
		viper.SetConfigFile(flagConfig)
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

	viper.SetEnvPrefix("ukf")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if flagConfig != "" && err != nil {
		output.Fatalf("Failed to read config from file '%s': %s", flagConfig, err.Error())
	}
}
