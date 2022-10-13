package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/spf13/cobra"
)

func loadbalancerTargetGroupTargetRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "target",
		Short: "sub-commands relating to targets",
	}

	// Child commands
	cmd.AddCommand(loadbalancerTargetGroupTargetListCmd(f))
	cmd.AddCommand(loadbalancerTargetGroupTargetShowCmd(f))
	cmd.AddCommand(loadbalancerTargetGroupTargetCreateCmd(f))
	cmd.AddCommand(loadbalancerTargetGroupTargetUpdateCmd(f))
	cmd.AddCommand(loadbalancerTargetGroupTargetDeleteCmd(f))

	return cmd
}

func loadbalancerTargetGroupTargetListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <targetgroup: id>",
		Short:   "Lists targets",
		Long:    "This command lists targets",
		Example: "ans loadbalancer targetgroup target list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing target group")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerTargetGroupTargetList),
	}
}

func loadbalancerTargetGroupTargetList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	targetGroupID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid target group ID")
	}

	targets, err := service.GetTargetGroupTargets(targetGroupID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving targets: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerTargetsProvider(targets))
}

func loadbalancerTargetGroupTargetShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <targetgroup: id> <target: id>...",
		Short:   "Shows a target",
		Long:    "This command shows one or more targets",
		Example: "ans loadbalancer targetgroup target show 123 345",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing target group")
			}
			if len(args) < 2 {
				return errors.New("Missing target")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerTargetGroupTargetShow),
	}
}

func loadbalancerTargetGroupTargetShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	targetGroupID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid target group ID")
	}

	var targets []loadbalancer.Target
	for _, arg := range args[1:] {

		targetID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid target ID [%s]", arg)
			continue
		}

		target, err := service.GetTargetGroupTarget(targetGroupID, targetID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving target [%d]: %s", targetID, err)
			continue
		}

		targets = append(targets, target)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerTargetsProvider(targets))
}

func loadbalancerTargetGroupTargetCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <targetgroup: id>",
		Short:   "Creates a target",
		Long:    "This command creates a target",
		Example: "ans loadbalancer targetgroup target create 123 --ip 1.2.3.4 --port 443",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing target group")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerTargetGroupTargetCreate),
	}

	cmd.Flags().String("name", "", "Name for target")
	cmd.Flags().String("ip", "", "IP address for target")
	cmd.MarkFlagRequired("ip")
	cmd.Flags().Int("port", 0, "Port number for target")
	cmd.MarkFlagRequired("port")
	cmd.Flags().Int("weight", 0, "Weight for target")
	cmd.Flags().Bool("backup", false, "Specifies whether target should be a backup")
	cmd.Flags().Int("check-interval", 0, "Check interval for target")
	cmd.Flags().Int("check-rise", 0, "Check rise for target")
	cmd.Flags().Int("check-fall", 0, "Check fall for target")
	cmd.Flags().Bool("disable-http2", false, "Specifies HTTP2 should be disabled for target")
	cmd.Flags().Bool("http2-only", false, "Specifies only HTTP2 should be enabled for target")
	cmd.Flags().Bool("active", true, "Specifies target should be active. Defaults to true")
	cmd.Flags().String("session-cookie-value", "", "Specifies value of session cookie for target")

	return cmd
}

func loadbalancerTargetGroupTargetCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	targetGroupID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid target group ID")
	}

	createRequest := loadbalancer.CreateTargetRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	ip, _ := cmd.Flags().GetString("ip")
	createRequest.IP = connection.IPAddress(ip)
	createRequest.Port, _ = cmd.Flags().GetInt("port")
	createRequest.Weight, _ = cmd.Flags().GetInt("weight")
	createRequest.Backup, _ = cmd.Flags().GetBool("backup")
	createRequest.CheckInterval, _ = cmd.Flags().GetInt("check-interval")
	createRequest.CheckSSL, _ = cmd.Flags().GetBool("check-ssl")
	createRequest.CheckRise, _ = cmd.Flags().GetInt("check-rise")
	createRequest.CheckFall, _ = cmd.Flags().GetInt("check-fall")
	createRequest.DisableHTTP2, _ = cmd.Flags().GetBool("disable-http2")
	createRequest.HTTP2Only, _ = cmd.Flags().GetBool("http2-only")
	createRequest.Active, _ = cmd.Flags().GetBool("active")
	createRequest.SessionCookieValue, _ = cmd.Flags().GetString("session-cookie-value")

	targetID, err := service.CreateTargetGroupTarget(targetGroupID, createRequest)
	if err != nil {
		return fmt.Errorf("Error creating target: %s", err)
	}

	target, err := service.GetTargetGroupTarget(targetGroupID, targetID)
	if err != nil {
		return fmt.Errorf("Error retrieving new target: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerTargetsProvider([]loadbalancer.Target{target}))
}

func loadbalancerTargetGroupTargetUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <targetgroup: id> <target: id>...",
		Short:   "Updates a target",
		Long:    "This command updates one or more targets",
		Example: "ans loadbalancer targetgroup target update 123 456 --port 443",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing target group")
			}
			if len(args) < 2 {
				return errors.New("Missing target")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerTargetGroupTargetUpdate),
	}

	cmd.Flags().String("name", "", "Name for target")
	cmd.Flags().String("ip", "", "IP address for target")
	cmd.Flags().Int("port", 0, "Port number for target")
	cmd.Flags().Int("weight", 0, "Weight for target")
	cmd.Flags().Bool("backup", false, "Specifies whether target should be a backup")
	cmd.Flags().Int("check-interval", 0, "Check interval for target")
	cmd.Flags().Int("check-rise", 0, "Check rise for target")
	cmd.Flags().Int("check-fall", 0, "Check fall for target")
	cmd.Flags().Bool("disable-http2", false, "Specifies HTTP2 should be disabled for target")
	cmd.Flags().Bool("http2-only", false, "Specifies only HTTP2 should be enabled for target")
	cmd.Flags().Bool("active", true, "Specifies target should be active")
	cmd.Flags().String("session-cookie-value", "", "Specifies value of session cookie for target")

	return cmd
}

func loadbalancerTargetGroupTargetUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	targetGroupID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid target group ID")
	}

	patchRequest := loadbalancer.PatchTargetRequest{}
	patchRequest.Name, _ = cmd.Flags().GetString("name")
	ip, _ := cmd.Flags().GetString("ip")
	patchRequest.IP = connection.IPAddress(ip)
	patchRequest.Port, _ = cmd.Flags().GetInt("port")
	patchRequest.Weight, _ = cmd.Flags().GetInt("weight")
	patchRequest.Backup = helper.GetBoolPtrFlagIfChanged(cmd, "backup")
	patchRequest.CheckInterval, _ = cmd.Flags().GetInt("check-interval")
	patchRequest.CheckSSL = helper.GetBoolPtrFlagIfChanged(cmd, "check-ssl")
	patchRequest.CheckRise, _ = cmd.Flags().GetInt("check-rise")
	patchRequest.CheckFall, _ = cmd.Flags().GetInt("check-fall")
	patchRequest.DisableHTTP2 = helper.GetBoolPtrFlagIfChanged(cmd, "disable-http2")
	patchRequest.HTTP2Only = helper.GetBoolPtrFlagIfChanged(cmd, "http2-only")
	patchRequest.Active = helper.GetBoolPtrFlagIfChanged(cmd, "active")
	patchRequest.SessionCookieValue, _ = cmd.Flags().GetString("session-cookie-value")

	var targets []loadbalancer.Target
	for _, arg := range args[1:] {
		targetID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid target ID [%s]", arg)
			continue
		}

		err = service.PatchTargetGroupTarget(targetGroupID, targetID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating target [%d]: %s", targetID, err)
			continue
		}

		target, err := service.GetTargetGroupTarget(targetGroupID, targetID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated target [%d]: %s", targetID, err)
			continue
		}

		targets = append(targets, target)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerTargetsProvider(targets))
}

func loadbalancerTargetGroupTargetDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <targetgroup: id> <target: id>...",
		Short:   "Removes a target",
		Long:    "This command removes one or more targets",
		Example: "ans loadbalancer targetgroup target delete 123 456",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing target group")
			}
			if len(args) < 2 {
				return errors.New("Missing target")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerTargetGroupTargetDelete),
	}
}

func loadbalancerTargetGroupTargetDelete(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	targetGroupID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid target group ID")
	}

	for _, arg := range args[1:] {
		targetID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid target ID [%s]", arg)
			continue
		}

		err = service.DeleteTargetGroupTarget(targetGroupID, targetID)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing target [%d]: %s", targetID, err)
			continue
		}
	}

	return nil
}
