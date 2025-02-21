package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/spf13/cobra"
)

func rawCmd(f connection.ConnectionFactory) *cobra.Command {
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
	cmd.Flags().StringArrayP("header", "H", []string{}, "Additional header for request")
	return cmd
}

type rawCommandOutput string

func (r *rawCommandOutput) Deserialize(resp *connection.APIResponse) error {
	defer resp.Response.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Response.Body)
	if err != nil {
		return err
	}

	*r = rawCommandOutput(string(bodyBytes))
	return nil
}

func (r *rawCommandOutput) Error() string {
	return string(*r)
}

type rawCommandData string

func (r *rawCommandData) Serialize() ([]byte, error) {
	return []byte(*r), nil
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
		data, _ := cmd.Flags().GetString("data")
		commandData := rawCommandData(data)
		req.Body = &commandData
	}

	if cmd.Flags().Changed("header") {
		req.Headers = http.Header{}
		headers, _ := cmd.Flags().GetStringArray("header")
		for _, header := range headers {
			headerParts := strings.SplitN(header, ":", 2)
			if len(headerParts) != 2 {
				return fmt.Errorf("invalid header format: %s", header)
			}

			req.Headers.Add(strings.TrimSpace(headerParts[0]), strings.TrimSpace(headerParts[1]))
		}
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

	fmt.Println(respBody)
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
