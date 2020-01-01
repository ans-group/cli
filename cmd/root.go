package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/blang/semver"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ukfast/cli/internal/pkg/build"
	"github.com/ukfast/cli/internal/pkg/output"
	apiclient "github.com/ukfast/sdk-go/pkg/client"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/logging"
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

	if err := rootCmd.Execute(); err != nil {
		output.Fatal(err.Error())
	}

	output.ExitWithErrorLevel()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&flagConfig, "config", "", "config file (default is $HOME/.ukfast.yaml)")
	rootCmd.PersistentFlags().StringVarP(&flagFormat, "format", "f", "", "output format {table, json, template, value, csv, list}")
	rootCmd.PersistentFlags().StringVar(&flagOutputTemplate, "outputtemplate", "", "output Go template (used with 'template' format), e.g. 'Name: {{ .Name }}'")
	rootCmd.PersistentFlags().StringVar(&flagSort, "sort", "", "output sorting, e.g. 'name', 'name:asc', 'name:desc'")
	rootCmd.PersistentFlags().StringSliceVar(&flagProperty, "property", []string{}, "property to output (used with several formats), can be repeated")
	rootCmd.PersistentFlags().StringArrayVar(&flagFilter, "filter", []string{}, "filter for list commands, can be repeated, e.g. 'property=somevalue', 'property:gt=3', 'property=valu*'")

	// Child commands
	rootCmd.AddCommand(updateCmd())

	// Child root commands
	rootCmd.AddCommand(completionRootCmd())
	rootCmd.AddCommand(safednsRootCmd())
	rootCmd.AddCommand(ecloudRootCmd())
	rootCmd.AddCommand(sslRootCmd())
	rootCmd.AddCommand(ddosxRootCmd())
	rootCmd.AddCommand(accountRootCmd())
	rootCmd.AddCommand(registrarRootCmd())
	rootCmd.AddCommand(pssRootCmd())
	rootCmd.AddCommand(storageRootCmd())
	rootCmd.AddCommand(loadtestRootCmd())

	appFilesystem = afero.NewOsFs()
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

func getClient() apiclient.Client {
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		output.Fatal("UKF_API_KEY not set")
	}

	apiTimeout := viper.GetInt("api_timeout_seconds")
	apiURI := viper.GetString("api_uri")
	apiInsecure := viper.GetBool("api_insecure")
	apiHeaders := viper.GetStringMapString("api_headers")
	apiDebug := viper.GetBool("api_debug")

	conn := connection.NewAPIConnection(&connection.APIKeyCredentials{APIKey: apiKey})
	conn.UserAgent = "ukfast-cli"
	if apiURI != "" {
		conn.APIURI = apiURI
	}
	if apiTimeout > 0 {
		conn.HTTPClient.Timeout = (time.Duration(apiTimeout) * time.Second)
	}
	if apiInsecure {
		conn.HTTPClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	if apiHeaders != nil {
		conn.Headers = http.Header{}
		for headerKey, headerValue := range apiHeaders {
			conn.Headers.Add(headerKey, headerValue)
		}
	}

	if apiDebug {
		logging.SetLogger(&output.DebugLogger{})
	}

	return apiclient.NewClient(conn)
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
