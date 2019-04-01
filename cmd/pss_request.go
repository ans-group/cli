package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/pss"
)

func pssRequestRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request",
		Short: "sub-commands relating to requests",
	}

	// Child commands
	cmd.AddCommand(pssRequestListCmd())
	cmd.AddCommand(pssRequestShowCmd())

	// Child root commands
	cmd.AddCommand(pssRequestReplyRootCmd())

	return cmd
}

func pssRequestListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists requests",
		Long:    "This command lists requests",
		Example: "ukfast pss request list",
		Run: func(cmd *cobra.Command, args []string) {
			pssRequestList(getClient().PSSService(), cmd, args)
		},
	}
}

func pssRequestList(service pss.PSSService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	requests, err := service.GetRequests(params)
	if err != nil {
		output.Fatalf("Error retrieving requests: %s", err)
		return
	}

	outputPSSRequests(requests)
}

func pssRequestShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <request: id>...",
		Short:   "Shows a request",
		Long:    "This command shows one or more requests",
		Example: "ukfast pss request show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing request")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			pssRequestShow(getClient().PSSService(), cmd, args)
		},
	}
}

func pssRequestShow(service pss.PSSService, cmd *cobra.Command, args []string) {
	var requests []pss.Request
	for _, arg := range args {
		requestID, err := strconv.Atoi(arg)
		if err != nil {
			OutputWithErrorLevelf("Invalid request ID [%s]", arg)
			continue
		}

		request, err := service.GetRequest(requestID)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving request [%s]: %s", arg, err)
			continue
		}

		requests = append(requests, request)
	}

	outputPSSRequests(requests)
}
