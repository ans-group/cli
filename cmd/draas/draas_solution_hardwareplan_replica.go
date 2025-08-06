package draas

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/draas"
	"github.com/spf13/cobra"
)

func draasSolutionHardwarePlanReplicaRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replica",
		Short: "sub-commands relating to solution hardware plan replicas",
	}

	// Child commands
	cmd.AddCommand(draasSolutionHardwarePlanReplicaListCmd(f))

	return cmd
}

func draasSolutionHardwarePlanReplicaListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <solution: id> <hardwareplan: id>",
		Short:   "Lists solution harware plan replicas",
		Long:    "This command lists solution harware plan replicas",
		Example: "ans draas solution hardwareplan replica list 00000000-0000-0000-0000-000000000000 00000000-0000-0000-0000-000000000001",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing solution")
			}
			if len(args) < 2 {
				return errors.New("missing hardware plan")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return draasSolutionHardwarePlanReplicaList(c.DRaaSService(), cmd, args)
		},
	}
}

func draasSolutionHardwarePlanReplicaList(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	replicas, err := service.GetSolutionHardwarePlanReplicas(args[0], args[1], params)
	if err != nil {
		return fmt.Errorf("error retrieving solution hardware plan replicas: %s", err.Error())
	}

	return output.CommandOutput(cmd, ReplicaCollection(replicas))
}
