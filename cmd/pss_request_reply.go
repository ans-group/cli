package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/pss"
)

func pssRequestReplyRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reply",
		Short: "sub-commands relating to request replies",
	}

	// Child commands
	cmd.AddCommand(pssRequestReplyListCmd())

	// Child root commands

	return cmd
}

func pssRequestReplyListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list <request: id>",
		Short:   "Lists a request",
		Long:    "This command lists the replies for a request",
		Example: "ukfast pss request reply list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing request")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			pssRequestReplyList(getClient().PSSService(), cmd, args)
		},
	}
}

func pssRequestReplyList(service pss.PSSService, cmd *cobra.Command, args []string) {
	requestID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid request ID [%s]", args[0])
		return
	}

	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	replies, err := service.GetRequestConversation(requestID, params)
	if err != nil {
		output.Fatalf("Error retrieving request replies: %s", err)
		return
	}

	outputPSSReplies(replies)
}
