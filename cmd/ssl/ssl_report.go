package ssl

import (
	"errors"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ssl"
	"github.com/spf13/cobra"
)

func sslReportRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "sub-commands relating to reports",
	}

	// Child commands
	cmd.AddCommand(sslReportShowCmd(f))

	return cmd
}

func sslReportShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <report: id>...",
		Short:   "Shows a report",
		Long:    "This command shows one or more reports",
		Example: "ans ssl report show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return sslReportShow(c.SSLService(), cmd, args)
		},
	}
}

func sslReportShow(service ssl.SSLService, cmd *cobra.Command, args []string) error {
	var reports []ssl.Report
	for _, arg := range args {
		report, err := service.GetReport(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving SSL report for domain [%s]: %s", arg, err)
			continue
		}

		reports = append(reports, report)
	}

	return output.CommandOutput(cmd, ReportCollection(reports))
}
