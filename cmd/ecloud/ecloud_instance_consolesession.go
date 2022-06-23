package ecloud

import (
	"errors"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/pkg/browser"

	"github.com/spf13/cobra"
)

func ecloudInstanceConsoleRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "consolesession",
		Short: "sub-commands relating to instance Consoles",
	}

	// Child commands
	cmd.AddCommand(ecloudInstanceConsoleSessionCreateCmd(f))

	return cmd
}

func ecloudInstanceConsoleSessionCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <instance: id>",
		Short:   "Creates an instance console session",
		Long:    "This command creates one or more instance console sessions",
		Example: "ans ecloud instance consolesession create i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceConsoleSessionCreate),
	}

	cmd.Flags().Bool("browser", false, "Indicates session should be opened in default browser")

	return cmd
}

func ecloudInstanceConsoleSessionCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var sessions []ecloud.ConsoleSession
	for _, arg := range args {
		session, err := service.CreateInstanceConsoleSession(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error creating instance [%s] console session: %s", arg, err)
			continue
		}

		openBrowser, _ := cmd.Flags().GetBool("browser")
		if openBrowser {
			err = browser.OpenURL(session.URL)
			if err != nil {
				output.OutputWithErrorLevelf("Error opening console session in browser for instance [%s]: %s", arg, err)
			}
		}

		sessions = append(sessions, session)
	}

	return output.CommandOutput(cmd, OutputECloudConsoleSessionsProvider(sessions))
}
