package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/sdk-go/pkg/service/pss"
)

func pssReplyRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reply",
		Short: "sub-commands relating to replies",
	}

	// Child commands
	cmd.AddCommand(pssReplyShowCmd())

	return cmd
}

func pssReplyShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <reply: id>...",
		Short:   "Shows a reply",
		Long:    "This command shows one or more replies",
		Example: "ukfast pss reply show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing reply")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			pssReplyShow(getClient().PSSService(), cmd, args)
		},
	}
}

func pssReplyShow(service pss.PSSService, cmd *cobra.Command, args []string) {
	var replies []pss.Reply
	for _, arg := range args {
		reply, err := service.GetReply(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving reply [%s]: %s", arg, err)
			continue
		}

		replies = append(replies, reply)
	}

	outputPSSReplies(replies)
}
