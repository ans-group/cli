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
	cmd.AddCommand(pssRequestReplyCreateCmd())

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

func pssRequestReplyCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a reply",
		Long:    "This command creates a new reply",
		Example: "ukfast pss request reply create --description 'example' --author 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing request")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			pssRequestReplyCreate(getClient().PSSService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("description", "", "Specifies description for reply")
	cmd.MarkFlagRequired("description")
	cmd.Flags().Int("author", 0, "Specifies author ID for reply")
	cmd.MarkFlagRequired("author")

	return cmd
}

func pssRequestReplyCreate(service pss.PSSService, cmd *cobra.Command, args []string) {
	requestID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid request ID [%s]", args[0])
		return
	}

	createRequest := pss.CreateReplyRequest{}
	createRequest.Author.ID, _ = cmd.Flags().GetInt("author")
	createRequest.Description, _ = cmd.Flags().GetString("description")

	replyID, err := service.CreateRequestReply(requestID, createRequest)
	if err != nil {
		output.Fatalf("Error creating reply: %s", err)
		return
	}

	reply, err := service.GetReply(replyID)
	if err != nil {
		output.Fatalf("Error retrieving new reply: %s", err)
		return
	}

	outputPSSReplies([]pss.Reply{reply})
}
