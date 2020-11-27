package loadbalancer

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	flaghelper "github.com/ukfast/cli/internal/pkg/helper/flag"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func loadbalancerTargetRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "target",
		Short: "sub-commands relating to targets",
	}

	// Child commands
	cmd.AddCommand(loadbalancerTargetListCmd(f))
	cmd.AddCommand(loadbalancerTargetShowCmd(f))
	cmd.AddCommand(loadbalancerTargetCreateCmd(f))
	cmd.AddCommand(loadbalancerTargetUpdateCmd(f))
	cmd.AddCommand(loadbalancerTargetDeleteCmd(f))

	return cmd
}

func loadbalancerTargetListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists targets",
		Long:    "This command lists targets",
		Example: "ukfast loadbalancer target list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerTargetList(c.LoadBalancerService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Target name for filtering")

	return cmd
}

func loadbalancerTargetList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := flaghelper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	flaghelper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, flaghelper.NewStringFilterFlag("name", "name"))

	targets, err := service.GetTargets(params)
	if err != nil {
		return fmt.Errorf("Error retrieving targets: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerTargetsProvider(targets))
}

func loadbalancerTargetShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <target: id>...",
		Short:   "Shows a target",
		Long:    "This command shows one or more targets",
		Example: "ukfast loadbalancer target show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing target")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerTargetShow(c.LoadBalancerService(), cmd, args)
		},
	}
}

func loadbalancerTargetShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var targets []loadbalancer.Target
	for _, arg := range args {
		target, err := service.GetTarget(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving target [%s]: %s", arg, err)
			continue
		}

		targets = append(targets, target)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerTargetsProvider(targets))
}

func loadbalancerTargetCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <target: id>...",
		Short:   "Creates a target",
		Long:    "This command creates a target",
		Example: "ukfast loadbalancer target create",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerTargetCreate(c.LoadBalancerService(), cmd, args)
		},
	}

	cmd.Flags().String("target-group", "", "Group for target")
	cmd.Flags().String("ip", "", "IP address for target")
	cmd.Flags().Int("port", 0, "Port number for target")
	cmd.Flags().Int("weight", 0, "Weight of target")
	cmd.Flags().Bool("backup", false, "Specifies target should be backup")
	cmd.Flags().Int("check-interval", 0, "Interval between checks")
	cmd.Flags().Bool("check-ssl", false, "Specifies checks should be performed using SSL")
	cmd.Flags().Int("check-rise", 0, "Specifies rise value for checks")
	cmd.Flags().Int("check-fall", 0, "Specifies fall value for checks")
	cmd.Flags().Bool("disable-http2", false, "Specifies HTTP2 should be disabled for target")
	cmd.Flags().Bool("http2-only", false, "Specifies only HTTP2 should be enabled for target")

	return cmd
}

func loadbalancerTargetCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	req := loadbalancer.CreateTargetRequest{}
	req.TargetGroupID, _ = cmd.Flags().GetString("target-group")
	ip, _ := cmd.Flags().GetString("ip")
	req.IP = connection.IPAddress(ip)
	req.Port, _ = cmd.Flags().GetInt("port")
	req.Weight, _ = cmd.Flags().GetInt("weight")
	req.Backup, _ = cmd.Flags().GetBool("backup")
	req.CheckInterval, _ = cmd.Flags().GetInt("check-interval")
	req.CheckSSL, _ = cmd.Flags().GetBool("check-ssl")
	req.CheckRise, _ = cmd.Flags().GetInt("check-rise")
	req.CheckFall, _ = cmd.Flags().GetInt("check-fall")
	req.DisableHTTP2, _ = cmd.Flags().GetBool("disable-http2")
	req.HTTP2Only, _ = cmd.Flags().GetBool("http2-only")

	targetID, err := service.CreateTarget(req)
	if err != nil {
		return fmt.Errorf("Error creating target: %s", err)
	}

	target, err := service.GetTarget(targetID)
	if err != nil {
		return fmt.Errorf("Error retrieving new target: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerTargetsProvider([]loadbalancer.Target{target}))
}

func loadbalancerTargetUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <target: id>...",
		Short:   "Updates a target",
		Long:    "This command updates a target",
		Example: "ukfast loadbalancer target update",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing target")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerTargetUpdate(c.LoadBalancerService(), cmd, args)
		},
	}

	cmd.Flags().String("target-group", "", "Group for target")
	cmd.Flags().String("ip", "", "IP address for target")
	cmd.Flags().Int("port", 0, "Port number for target")
	cmd.Flags().Int("weight", 0, "Weight of target")
	cmd.Flags().Bool("backup", false, "Specifies target should be backup")
	cmd.Flags().Int("check-interval", 0, "Interval between checks")
	cmd.Flags().Bool("check-ssl", false, "Specifies checks should be performed using SSL")
	cmd.Flags().Int("check-rise", 0, "Specifies rise value for checks")
	cmd.Flags().Int("check-fall", 0, "Specifies fall value for checks")
	cmd.Flags().Bool("disable-http2", false, "Specifies HTTP2 should be disabled for target")
	cmd.Flags().Bool("http2-only", false, "Specifies only HTTP2 should be enabled for target")

	return cmd
}

func loadbalancerTargetUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var targets []loadbalancer.Target

	req := loadbalancer.PatchTargetRequest{}
	req.TargetGroupID, _ = cmd.Flags().GetString("target-group")
	ip, _ := cmd.Flags().GetString("ip")
	req.IP = connection.IPAddress(ip)
	req.Port, _ = cmd.Flags().GetInt("port")
	req.Weight, _ = cmd.Flags().GetInt("weight")
	req.Backup = flaghelper.GetChangedOrDefaultPtrBool(cmd, "backup")
	req.CheckInterval, _ = cmd.Flags().GetInt("check-interval")
	req.CheckSSL = flaghelper.GetChangedOrDefaultPtrBool(cmd, "check-ssl")
	req.CheckRise, _ = cmd.Flags().GetInt("check-rise")
	req.CheckFall, _ = cmd.Flags().GetInt("check-fall")
	req.DisableHTTP2 = flaghelper.GetChangedOrDefaultPtrBool(cmd, "disable-http2")
	req.HTTP2Only = flaghelper.GetChangedOrDefaultPtrBool(cmd, "http2-only")

	for _, arg := range args[1:] {
		err := service.PatchTarget(arg, req)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating target [%s]: %s", arg, err)
			continue
		}

		target, err := service.GetTarget(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated target [%s]: %s", arg, err)
			continue
		}

		targets = append(targets, target)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerTargetsProvider(targets))
}

func loadbalancerTargetDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <target: id>...",
		Short:   "Deletes a target",
		Long:    "This command deletes a target",
		Example: "ukfast loadbalancer target delete",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing target")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerTargetDelete(c.LoadBalancerService(), cmd, args)
		},
	}
}

func loadbalancerTargetDelete(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	for _, arg := range args[1:] {
		err := service.DeleteTarget(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing target [%s]: %s", arg, err)
			continue
		}
	}

	return nil
}
