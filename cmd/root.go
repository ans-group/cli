package cmd

import (
	"errors"
	"fmt"

	"github.com/blang/semver"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
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
var flagSort string
var flagProperty []string
var flagFilter []string
var appFilesystem afero.Fs
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
	rootCmd.PersistentFlags().StringVar(&flagConfig, "config", "", "config file (default is $HOME/.ukfast.yaml)")
	rootCmd.PersistentFlags().StringVarP(&flagFormat, "format", "f", "", "output format {table, json, template, value, csv, list}")
	rootCmd.PersistentFlags().StringVar(&flagOutputTemplate, "outputtemplate", "", "output Go template (used with 'template' format), e.g. 'Name: {{ .Name }}'")
	rootCmd.PersistentFlags().StringVar(&flagSort, "sort", "", "output sorting, e.g. 'name', 'name:asc', 'name:desc'")
	rootCmd.PersistentFlags().StringSliceVar(&flagProperty, "property", []string{}, "property to output (used with several formats), can be repeated")
	rootCmd.PersistentFlags().StringArrayVar(&flagFilter, "filter", []string{}, "filter for list commands, can be repeated, e.g. 'property=somevalue', 'property:gt=3', 'property=valu*'")

	initConfig()

	appFilesystem = afero.NewOsFs()

	clientFactory, err := getClientFactory()
	if err != nil {
		output.Fatal(err.Error())
	}

	// Child commands
	rootCmd.AddCommand(updateCmd())

	// Child root commands
	rootCmd.AddCommand(CompletionRootCmd())
	rootCmd.AddCommand(accountcmd.AccountRootCmd(clientFactory))
	rootCmd.AddCommand(ddosxcmd.DDoSXRootCmd(clientFactory, appFilesystem))
	rootCmd.AddCommand(ecloudcmd.ECloudRootCmd(clientFactory))
	rootCmd.AddCommand(loadtestcmd.LoadTestRootCmd(clientFactory))
	rootCmd.AddCommand(psscmd.PSSRootCmd(clientFactory, appFilesystem))
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
	}

	viper.SetEnvPrefix("ukf")
	viper.SetDefault("api_timeout_seconds", 90)
	viper.SetDefault("command_wait_timeout_seconds", 1200)
	viper.SetDefault("command_wait_sleep_seconds", 5)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}

func getClientFactory() (factory.ClientFactory, error) {
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		return nil, errors.New("UKF_API_KEY not set")
	}

	return factory.NewUKFastClientFactory(
		factory.WithAPIKey(apiKey),
		factory.WithTimeout(viper.GetInt("api_timeout_seconds")),
		factory.WithURI(viper.GetString("api_uri")),
		factory.WithInsecure(viper.GetBool("api_insecure")),
		factory.WithHeaders(viper.GetStringMapString("api_headers")),
		factory.WithDebug(viper.GetBool("api_debug")),
	), nil
}

func updateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Updates CLI to latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			currentVersion, err := semver.ParseTolerant(appVersion)
			if err != nil {
				return fmt.Errorf("Unable to parse version: %s", err.Error())
			}
			newRelease, err := selfupdate.UpdateSelf(currentVersion, "ukfast/cli")
			if err != nil {
				return fmt.Errorf("Failed to update UKFast CLI: %s", err.Error())
			}

			if currentVersion.Equals(newRelease.Version) {
				fmt.Printf("UKFast CLI already at latest version (%s)\n", appVersion)
			} else {
				fmt.Printf("UKFast CLI updated to version v%s successfully\n", newRelease.Version)
				fmt.Println("Release notes:\n", newRelease.ReleaseNotes)
			}
			return nil
		},
	}
}
