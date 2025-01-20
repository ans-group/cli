package ddosx

import (
	"errors"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
)

func ddosxSSLContentRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "content",
		Short: "sub-commands relating to content",
	}

	// Child commands
	cmd.AddCommand(ddosxSSLContentShowCmd(f))

	return cmd
}

func ddosxSSLContentShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <ssl: id>...",
		Short:   "Shows a ssl's content",
		Long:    "This command shows one or more ssl's content",
		Example: "ans ddosx ssl content show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ssl")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxSSLContentShow(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxSSLContentShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var sslContents []ddosx.SSLContent
	for _, arg := range args {
		sslContent, err := service.GetSSLContent(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving ssl [%s]: %s", arg, err)
			continue
		}

		sslContents = append(sslContents, sslContent)
	}

	return output.CommandOutput(cmd, SSLContentCollection(sslContents))
}
