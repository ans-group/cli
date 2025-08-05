package draas

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/sdk-go/pkg/service/draas"
	"github.com/spf13/cobra"
)

func draasSolutionReplicaIOPSRootCmd(f factory.ClientFactory) *cobra.Command { //nolint:unused
	cmd := &cobra.Command{
		Use:   "iops",
		Short: "sub-commands relating to solution replica IOPS tiers",
	}

	// Child commands

	return cmd
}

func draasSolutionReplicaIOPSUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <solution: id> <replica: id>...",
		Short:   "Updates the IOPS for a replica",
		Long:    "This command updates the IOPS for one or more replicas",
		Example: "ans draas solution update 00000000-0000-0000-0000-000000000000 --name test",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing solution")
			}
			if len(args) < 2 {
				return errors.New("missing replica")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return draasSolutionReplicaIOPSUpdate(c.DRaaSService(), cmd, args)
		},
	}

	cmd.Flags().String("iops-tier", "", "IOPS tier ID")
	_ = cmd.MarkFlagRequired("iops-tier")

	return cmd
}

func draasSolutionReplicaIOPSUpdate(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	req := draas.UpdateReplicaIOPSRequest{}

	iopsTierID, _ := cmd.Flags().GetString("iops-tier")
	req.IOPSTierID = iopsTierID

	for _, arg := range args[1:] {
		err := service.UpdateSolutionReplicaIOPS(args[0], arg, req)
		if err != nil {
			return fmt.Errorf("error updating replica [%s]: %s", arg, err.Error())
		}
	}

	return nil
}
