package cloudflare

import (
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/cloudflare"
	"github.com/spf13/cobra"
)

func cloudflareTotalSpendRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "totalspend",
		Short: "sub-commands relating to total spend",
	}

	// Child commands
	cmd.AddCommand(cloudflareTotalSpendShowCmd(f))

	return cmd
}

func cloudflareTotalSpendShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Shows total spend",
		Long:    "This command shows total spend",
		Example: "ans cloudflare totalspend show",
		RunE:    cloudflareCobraRunEFunc(f, cloudflareTotalSpendShow),
	}
}

func cloudflareTotalSpendShow(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	spend, err := service.GetTotalSpendMonthToDate()
	if err != nil {
		return fmt.Errorf("Error retrieving total spend: %s", err)
	}

	return output.CommandOutput(cmd, TotalSpendCollection([]cloudflare.TotalSpend{spend}))
}
