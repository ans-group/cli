package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxSSLContentRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "content",
		Short: "sub-commands relating to content",
	}

	// Child commands
	cmd.AddCommand(ddosxSSLContentShowCmd())

	return cmd
}

func ddosxSSLContentShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <ssl: id>...",
		Short:   "Shows a ssl's content",
		Long:    "This command shows one or more ssl's content",
		Example: "ukfast ddosx ssl content show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ssl")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxSSLContentShow(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxSSLContentShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var sslContents []ddosx.SSLContent
	for _, arg := range args {
		sslContent, err := service.GetSSLContent(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving ssl [%s]: %s", arg, err)
			continue
		}

		sslContents = append(sslContents, sslContent)
	}

	outputDDoSXSSLContents(sslContents)
}
