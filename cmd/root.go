// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ukfast/cli/internal/pkg/build"
	"github.com/ukfast/cli/internal/pkg/helper"
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
var errorLevel int
var appFilesystem afero.Fs

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ukfast",
	Short:   "Utility for manipulating UKFast services",
	Version: "UNKNOWN",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(build build.BuildInfo) {
	rootCmd.Version = build.String()

	if err := rootCmd.Execute(); err != nil {
		output.Fatal(err.Error())
	}

	if errorLevel > 0 {
		os.Exit(errorLevel)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&flagConfig, "config", "", "config file (default is $HOME/.ukfast.yaml)")
	rootCmd.PersistentFlags().StringVarP(&flagFormat, "format", "f", "", "output format {table, json, template, value, csv}")
	rootCmd.PersistentFlags().StringVar(&flagOutputTemplate, "outputtemplate", "", "output Go template (used with 'template' format), e.g. 'Name: {{ .Name }}'")
	rootCmd.PersistentFlags().StringVar(&flagSort, "sort", "", "output sorting, e.g. 'name', 'name:asc', 'name:desc'")
	rootCmd.PersistentFlags().StringSliceVar(&flagProperty, "property", []string{}, "property to output (used with several formats), can be repeated")
	rootCmd.PersistentFlags().StringArrayVar(&flagFilter, "filter", []string{}, "filter for list commands, can be repeated, e.g. 'property=somevalue', 'property:gt=3', 'property=valu*'")

	// Child root commands
	rootCmd.AddCommand(safednsRootCmd())
	rootCmd.AddCommand(ecloudRootCmd())
	rootCmd.AddCommand(sslRootCmd())
	rootCmd.AddCommand(ddosxRootCmd())
	rootCmd.AddCommand(accountRootCmd())

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

// OutputWithCustomErrorLevel is a wrapper for OutputError, which sets global
// var errorLevel with provided level
func OutputWithCustomErrorLevel(level int, str string) {
	output.Error(str)
	errorLevel = level
}

// OutputWithCustomErrorLevelf is a wrapper for OutputWithCustomErrorLevel, which sets global
// var errorLevel with provided level
func OutputWithCustomErrorLevelf(level int, format string, a ...interface{}) {
	OutputWithCustomErrorLevel(level, fmt.Sprintf(format, a...))
}

// OutputWithErrorLevelf is a wrapper for OutputWithCustomErrorLevelf, which sets global
// var errorLevel to 1
func OutputWithErrorLevelf(format string, a ...interface{}) {
	OutputWithCustomErrorLevelf(1, format, a...)
}

// OutputWithErrorLevel is a wrapper for OutputWithCustomErrorLevel, which sets global
// var errorLevel to 1
func OutputWithErrorLevel(str string) {
	OutputWithCustomErrorLevel(1, str)
}

type OutputDataProvider interface {
	GetData() interface{}
	GetFieldData() ([]*output.OrderedFields, error)
}

// Output calls the relevant OutputProvider data retrieval methods for given value
// in global variable 'format'
func Output(out OutputDataProvider) error {
	format := flagFormat
	if format == "" {
		format = "table"
	}

	switch format {
	case "json":
		return output.JSON(out.GetData())
	case "template":
		return output.Template(flagOutputTemplate, out.GetData())
	case "value":
		d, err := out.GetFieldData()
		if err != nil {
			return err
		}
		return output.Value(flagProperty, d)
	case "csv":
		d, err := out.GetFieldData()
		if err != nil {
			return err
		}
		return output.CSV(flagProperty, d)
	default:
		output.Errorf("Invalid output format [%s], defaulting to 'table'", format)
		fallthrough
	case "table":
		d, err := out.GetFieldData()
		if err != nil {
			return err
		}
		return output.Table(flagProperty, d)
	}
}

type APIListParameters struct {
	Filtering []connection.APIRequestFiltering
	Sorting   connection.APIRequestSorting
}

func GetAPIRequestParametersFromFlags() (connection.APIRequestParameters, error) {
	filtering, err := helper.GetFilteringArrayFromStringArrayFlag(flagFilter)
	if err != nil {
		return connection.APIRequestParameters{}, err
	}

	return connection.APIRequestParameters{
		Sorting:   helper.GetSortingFromStringFlag(flagSort),
		Filtering: filtering,
		Pagination: connection.APIRequestPagination{
			PerPage: viper.GetInt("api_pagination_perpage"),
		},
	}, nil
}

type WaitFunc func() (finished bool, err error)

func WaitForCommand(f WaitFunc) error {
	waitTimeout := viper.GetInt("command_wait_timeout_seconds")
	if waitTimeout < 1 {
		return errors.New("Invalid command_wait_timeout_seconds")
	}
	sleepTimeout := viper.GetInt("command_wait_sleep_seconds")
	if sleepTimeout < 1 {
		return errors.New("Invalid command_wait_sleep_seconds")
	}

	timeStart := time.Now()

	for {
		if time.Since(timeStart).Seconds() > float64(waitTimeout) {
			return errors.New("Timed out waiting for command")
		}

		finished, err := f()
		if err != nil {
			return fmt.Errorf("Error waiting for command: %s", err)
		}
		if finished {
			break
		}

		time.Sleep(time.Duration(sleepTimeout) * time.Second)
	}

	return nil
}
