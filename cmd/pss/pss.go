package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func PSSRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pss",
		Short: "Commands relating to PSS service",
	}

	// Child root commands
	cmd.AddCommand(pssRequestRootCmd(f))
	cmd.AddCommand(pssReplyRootCmd(f, fs))
	cmd.AddCommand(pssIncidentRootCmd(f))
	cmd.AddCommand(pssChangeRootCmd(f))
	cmd.AddCommand(pssCaseRootCmd(f))
	cmd.AddCommand(pssProblemRootCmd(f))
	cmd.AddCommand(pssSupportedServiceRootCmd(f))

	return cmd
}

type pssServiceCobraRunEFunc func(service pss.PSSService, cmd *cobra.Command, args []string) error

func pssCobraRunEFunc(f factory.ClientFactory, rf pssServiceCobraRunEFunc) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		c, err := f.NewClient()
		if err != nil {
			return err
		}

		return rf(c.PSSService(), cmd, args)
	}
}
