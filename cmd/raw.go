package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/spf13/cobra"
)

func rawCmd(f factory.ConnectionFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raw",
		Short: "Executes raw commands against API",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing resource/uri")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewConnection()
			if err != nil {
				return err
			}

			return raw(c, cmd, args)
		},
	}

	cmd.Flags().String("method", "GET", "Method for request")
	cmd.Flags().StringP("request", "X", "GET", "Method for request (curl alias for 'method')")
	cmd.Flags().StringP("data", "d", "", "Data for request")
	return cmd
}

type rawCommandOutput string

func (r *rawCommandOutput) Deserializer() func(resp *connection.APIResponse, out interface{}) error {
	return func(resp *connection.APIResponse, out interface{}) error {
		defer resp.Response.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Response.Body)
		if err != nil {
			return err
		}

		outRef, _ := out.(*rawCommandOutput)
		*outRef = rawCommandOutput(string(bodyBytes))
		return nil
	}
}

func raw(c connection.Connection, cmd *cobra.Command, args []string) error {
	req := connection.APIRequest{
		Resource: args[0],
	}

	method, methodSet := getFlagStringValue(cmd, "method", "request")
	if !methodSet {
		method = "GET"
	}

	req.Method = strings.ToUpper(method)

	if method == "GET" {
		params, err := helper.GetAPIRequestParametersFromFlags(cmd)
		if err != nil {
			return err
		}

		req.Parameters = params
	}

	if cmd.Flags().Changed("data") {
		req.Body, _ = cmd.Flags().GetString("data")
	}

	resp, err := c.Invoke(req)
	if err != nil {
		return err
	}

	respBody := rawCommandOutput("")

	err = resp.HandleResponse(&respBody)
	if err != nil {
		return err
	}

	fmt.Print(respBody)
	return nil
}

func getFlagStringValue(cmd *cobra.Command, flagName string, flagAliases ...string) (string, bool) {
	if cmd.Flags().Changed(flagName) {
		value, _ := cmd.Flags().GetString(flagName)
		return value, true
	}

	for _, flagAlias := range flagAliases {
		if cmd.Flags().Changed(flagAlias) {
			value, _ := cmd.Flags().GetString(flagAlias)
			return value, true
		}
	}

	return "", false
}
