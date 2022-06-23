package pss

import (
	"errors"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func pssReplyRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reply",
		Short: "sub-commands relating to replies",
	}

	// Child commands
	cmd.AddCommand(pssReplyShowCmd(f))

	// Child root commands
	cmd.AddCommand(pssReplyAttachmentRootCmd(f, fs))

	return cmd
}

func pssReplyShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <reply: id>...",
		Short:   "Shows a reply",
		Long:    "This command shows one or more replies",
		Example: "ans pss reply show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing reply")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return pssReplyShow(c.PSSService(), cmd, args)
		},
	}
}

func pssReplyShow(service pss.PSSService, cmd *cobra.Command, args []string) error {
	var replies []pss.Reply
	for _, arg := range args {
		reply, err := service.GetReply(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving reply [%s]: %s", arg, err)
			continue
		}

		replies = append(replies, reply)
	}

	return output.CommandOutput(cmd, OutputPSSRepliesProvider(replies))
}
