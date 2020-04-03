package pss

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/input"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/pss"
)

func pssRequestReplyRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reply",
		Short: "sub-commands relating to request replies",
	}

	// Child commands
	cmd.AddCommand(pssRequestReplyListCmd(f))
	cmd.AddCommand(pssRequestReplyCreateCmd(f))

	// Child root commands

	return cmd
}

func pssRequestReplyListCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return pssRequestReplyList(f.NewClient().PSSService(), cmd, args)
		},
	}
}

func pssRequestReplyList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	requestID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid request ID [%s]", args[0])
	}

	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	replies, err := service.GetRequestConversation(requestID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving request replies: %s", err)
	}

	return output.CommandOutput(cmd, OutputPSSRepliesProvider(replies))
}

func pssRequestReplyCreateCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return pssRequestReplyCreate(f.NewClient().PSSService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("description", "", "Specifies description for reply")
	cmd.Flags().Int("author", 0, "Specifies author ID for reply")
	cmd.MarkFlagRequired("author")

	return cmd
}

func pssRequestReplyCreate(service pss.PSSService, cmd *cobra.Command, args []string) error {
	requestID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid request ID [%s]", args[0])
	}

	createRequest := pss.CreateReplyRequest{}
	createRequest.Author.ID, _ = cmd.Flags().GetInt("author")

	if cmd.Flags().Changed("description") {
		createRequest.Description, _ = cmd.Flags().GetString("description")
	} else {
		createRequest.Description, err = input.ReadInput("description")
		if err != nil {
			return err
		}
	}

	replyID, err := service.CreateRequestReply(requestID, createRequest)
	if err != nil {
		return fmt.Errorf("Error creating reply: %s", err)
	}

	reply, err := service.GetReply(replyID)
	if err != nil {
		return fmt.Errorf("Error retrieving new reply: %s", err)
	}

	return output.CommandOutput(cmd, OutputPSSRepliesProvider([]pss.Reply{reply}))
}
