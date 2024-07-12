package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/spf13/cobra"
)

func pssCaseOptionRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "caseoption",
		Short: "sub-commands relating to case options",
	}

	// Child root commands
	cmd.AddCommand(pssCaseOptionChangeImpactRootCmd(f))
	cmd.AddCommand(pssCaseOptionChangeRiskRootCmd(f))
	cmd.AddCommand(pssCaseOptionIncidentImpactRootCmd(f))
	cmd.AddCommand(pssCaseOptionIncidentTypeRootCmd(f))

	return cmd
}
