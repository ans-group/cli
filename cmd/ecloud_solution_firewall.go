package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudSolutionFirewallRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "firewall",
		Short: "sub-commands relating to solution firewalls",
	}

	// Child commands
	cmd.AddCommand(ecloudSolutionFirewallListCmd())

	return cmd
}

func ecloudSolutionFirewallListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists solution firewalls",
		Long:    "This command lists solution firewalls",
		Example: "ukfast ecloud solution firewall list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing solution")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudSolutionFirewallList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudSolutionFirewallList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	solutionID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid solution ID [%s]", args[0])
		return
	}

	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	firewalls, err := service.GetSolutionFirewalls(solutionID, params)
	if err != nil {
		output.Fatalf("Error retrieving solution firewalls: %s", err)
		return
	}

	outputECloudFirewalls(firewalls)
}
