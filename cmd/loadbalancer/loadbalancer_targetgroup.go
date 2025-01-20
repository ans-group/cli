package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/spf13/cobra"
)

func loadbalancerTargetGroupRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "targetgroup",
		Short: "sub-commands relating to target groups",
	}

	// Child commands
	cmd.AddCommand(loadbalancerTargetGroupListCmd(f))
	cmd.AddCommand(loadbalancerTargetGroupShowCmd(f))
	cmd.AddCommand(loadbalancerTargetGroupCreateCmd(f))
	cmd.AddCommand(loadbalancerTargetGroupUpdateCmd(f))
	cmd.AddCommand(loadbalancerTargetGroupDeleteCmd(f))

	// Child root commands
	cmd.AddCommand(loadbalancerTargetGroupACLRootCmd(f))
	cmd.AddCommand(loadbalancerTargetGroupTargetRootCmd(f))

	return cmd
}

func loadbalancerTargetGroupListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists target groups",
		Long:    "This command lists target groups",
		Example: "ans loadbalancer targetgroup list",
		RunE:    loadbalancerCobraRunEFunc(f, loadbalancerTargetGroupList),
	}
}

func loadbalancerTargetGroupList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	groups, err := service.GetTargetGroups(params)
	if err != nil {
		return fmt.Errorf("Error retrieving target groups: %s", err)
	}

	return output.CommandOutput(cmd, TargetGroupCollection(groups))
}

func loadbalancerTargetGroupShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <targetgroup: id>...",
		Short:   "Shows a target group",
		Long:    "This command shows one or more target groups",
		Example: "ans loadbalancer targetgroup show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing target group")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerTargetGroupShow),
	}
}

func loadbalancerTargetGroupShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var groups []loadbalancer.TargetGroup
	for _, arg := range args {
		groupID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid target group ID [%s]", arg)
			continue
		}

		group, err := service.GetTargetGroup(groupID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving target group [%d]: %s", groupID, err)
			continue
		}

		groups = append(groups, group)
	}

	return output.CommandOutput(cmd, TargetGroupCollection(groups))
}

func loadbalancerTargetGroupCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <targetgroup: id>",
		Short:   "Creates a target group",
		Long:    "This command creates a target group",
		Example: "ans loadbalancer targetgroup create --cluster 123 --name \"test-targetgroup\" --balance roundrobin --mode http",
		RunE:    loadbalancerCobraRunEFunc(f, loadbalancerTargetGroupCreate),
	}

	cmd.Flags().String("name", "", "Name of target group")
	cmd.MarkFlagRequired("name")
	cmd.Flags().Int("cluster", 0, "ID of cluster")
	cmd.MarkFlagRequired("cluster")
	cmd.Flags().String("balance", "", "Balance configuration for target group")
	cmd.MarkFlagRequired("balance")
	cmd.Flags().String("mode", "", "Specifies mode for target group")
	cmd.MarkFlagRequired("mode")
	cmd.Flags().Bool("close", false, "Specifies close should be enabled for target group")
	cmd.Flags().Bool("sticky", false, "Specifies sticky should be enabled for target group")
	cmd.Flags().String("cookie-opts", "", "Specifies cookie options for target group")
	cmd.Flags().String("source", "", "Specifies source for target group")
	cmd.Flags().Int("timeouts-connect", 0, "Specifies connect timeout for target group")
	cmd.Flags().Int("timeouts-server", 0, "Specifies server timeout for target group")
	cmd.Flags().Int("timeouts-http-request", 0, "Specifies HTTP request timeout for target group")
	cmd.Flags().Int("timeouts-check", 0, "Specifies check timeout for target group")
	cmd.Flags().Int("timeouts-tunnel", 0, "Specifies tunnel timeout for target group")
	cmd.Flags().String("custom-options", "", "Specifies custom options for target group")
	cmd.Flags().String("monitor-url", "", "Specifies monitor URL for target group")
	cmd.Flags().String("monitor-method", "", "Specifies monitor method for target group")
	cmd.Flags().String("monitor-host", "", "Specifies monitor host for target group")
	cmd.Flags().String("monitor-http-version", "", "Specifies monitor HTTP version for target group")
	cmd.Flags().String("monitor-expect", "", "Specifies monitor expected string for target group")
	cmd.Flags().String("monitor-expect-string", "", "Specifies monitor expected string for target group")
	cmd.Flags().Bool("monitor-expect-string-regex", false, "Specifies provided monitor expected string is a regular expression")
	cmd.Flags().Bool("monitor-tcp-monitoring", false, "Specifies monitor should use TCP for target group")
	cmd.Flags().Int("check-port", 0, "Specifies check port for target group")
	cmd.Flags().Bool("send-proxy", false, "Specifies proxy protocol should be used for target group")
	cmd.Flags().Bool("send-proxyv2", false, "Specifies proxy protocol V2 should be used for target group")
	cmd.Flags().Bool("ssl", false, "Specifies SSL should be used for target group")
	cmd.Flags().Bool("ssl-verify", false, "Specifies SSL verification should be used for target group")
	cmd.Flags().Bool("sni", false, "Specifies SNI should be used for target group")

	return cmd
}

func loadbalancerTargetGroupCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	createRequest := loadbalancer.CreateTargetGroupRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.ClusterID, _ = cmd.Flags().GetInt("cluster")
	createRequest.Close, _ = cmd.Flags().GetBool("close")
	createRequest.Sticky, _ = cmd.Flags().GetBool("sticky")
	createRequest.CookieOpts, _ = cmd.Flags().GetString("cookie-opts")
	createRequest.Source, _ = cmd.Flags().GetString("source")
	createRequest.TimeoutsConnect, _ = cmd.Flags().GetInt("timeouts-connect")
	createRequest.TimeoutsServer, _ = cmd.Flags().GetInt("timeouts-server")
	createRequest.TimeoutsHTTPRequest, _ = cmd.Flags().GetInt("timeouts-http-request")
	createRequest.TimeoutsCheck, _ = cmd.Flags().GetInt("timeouts-check")
	createRequest.TimeoutsTunnel, _ = cmd.Flags().GetInt("timeouts-tunnel")
	createRequest.CustomOptions, _ = cmd.Flags().GetString("custom-options")
	createRequest.MonitorURL, _ = cmd.Flags().GetString("monitor-url")
	createRequest.MonitorHost, _ = cmd.Flags().GetString("monitor-host")
	createRequest.MonitorHTTPVersion, _ = cmd.Flags().GetString("monitor-http-version")
	createRequest.MonitorExpect, _ = cmd.Flags().GetString("monitor-expect")
	createRequest.MonitorExpectString, _ = cmd.Flags().GetString("monitor-expect-string")
	createRequest.MonitorExpectStringRegex, _ = cmd.Flags().GetBool("monitor-expect-string-regex")
	createRequest.MonitorTCPMonitoring, _ = cmd.Flags().GetBool("monitor-tcp-monitoring")
	createRequest.CheckPort, _ = cmd.Flags().GetInt("check-port")
	createRequest.SendProxy, _ = cmd.Flags().GetBool("send-proxy")
	createRequest.SendProxyV2, _ = cmd.Flags().GetBool("send-proxyv2")
	createRequest.SSL, _ = cmd.Flags().GetBool("ssl")
	createRequest.SSLVerify, _ = cmd.Flags().GetBool("ssl-verify")
	createRequest.SNI, _ = cmd.Flags().GetBool("sni")

	balance, _ := cmd.Flags().GetString("balance")
	parsedBalance, err := loadbalancer.TargetGroupBalanceEnum.Parse(balance)
	if err != nil {
		return clierrors.NewErrInvalidFlagValue("balance", balance, err)
	}
	createRequest.Balance = parsedBalance

	mode, _ := cmd.Flags().GetString("mode")
	parsedMode, err := loadbalancer.ModeEnum.Parse(mode)
	if err != nil {
		return clierrors.NewErrInvalidFlagValue("mode", mode, err)
	}
	createRequest.Mode = parsedMode

	if cmd.Flags().Changed("monitor-method") {
		monitorMethod, _ := cmd.Flags().GetString("monitor-method")
		parsedMonitorMethod, err := loadbalancer.TargetGroupMonitorMethodEnum.Parse(monitorMethod)
		if err != nil {
			return clierrors.NewErrInvalidFlagValue("monitor-method", monitorMethod, err)
		}
		createRequest.MonitorMethod = parsedMonitorMethod
	}

	groupID, err := service.CreateTargetGroup(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating target group: %s", err)
	}

	group, err := service.GetTargetGroup(groupID)
	if err != nil {
		return fmt.Errorf("Error retrieving new target group: %s", err)
	}

	return output.CommandOutput(cmd, TargetGroupCollection([]loadbalancer.TargetGroup{group}))
}

func loadbalancerTargetGroupUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <targetgroup: id>...",
		Short:   "Updates a target group",
		Long:    "This command updates one or more target groups",
		Example: "ans loadbalancer targetgroup update 123 --name mytargetgroup",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing target group")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerTargetGroupUpdate),
	}
	cmd.Flags().String("name", "", "Name of target group")
	cmd.Flags().String("balance", "", "Balance configuration for target group")
	cmd.Flags().String("mode", "", "Specifies mode for target group")
	cmd.Flags().Bool("close", false, "Specifies close should be enabled for target group")
	cmd.Flags().Bool("sticky", false, "Specifies sticky should be enabled for target group")
	cmd.Flags().String("cookie-opts", "", "Specifies cookie options for target group")
	cmd.Flags().String("source", "", "Specifies source for target group")
	cmd.Flags().Int("timeouts-connect", 0, "Specifies connect timeout for target group")
	cmd.Flags().Int("timeouts-server", 0, "Specifies server timeout for target group")
	cmd.Flags().Int("timeouts-http-request", 0, "Specifies HTTP request timeout for target group")
	cmd.Flags().Int("timeouts-check", 0, "Specifies check timeout for target group")
	cmd.Flags().Int("timeouts-tunnel", 0, "Specifies tunnel timeout for target group")
	cmd.Flags().String("custom-options", "", "Specifies custom options for target group")
	cmd.Flags().String("monitor-url", "", "Specifies monitor URL for target group")
	cmd.Flags().String("monitor-method", "", "Specifies monitor method for target group")
	cmd.Flags().String("monitor-host", "", "Specifies monitor host for target group")
	cmd.Flags().String("monitor-http-version", "", "Specifies monitor HTTP version for target group")
	cmd.Flags().String("monitor-expect", "", "Specifies monitor expected string for target group")
	cmd.Flags().String("monitor-expect-string", "", "Specifies monitor expected string for target group")
	cmd.Flags().Bool("monitor-expect-string-regex", false, "Specifies provided monitor expected string is a regular expression")
	cmd.Flags().Bool("monitor-tcp-monitoring", false, "Specifies monitor should use TCP for target group")
	cmd.Flags().Int("check-port", 0, "Specifies check port for target group")
	cmd.Flags().Bool("send-proxy", false, "Specifies proxy protocol should be used for target group")
	cmd.Flags().Bool("send-proxyv2", false, "Specifies proxy protocol V2 should be used for target group")
	cmd.Flags().Bool("ssl", false, "Specifies SSL should be used for target group")
	cmd.Flags().Bool("ssl-verify", false, "Specifies SSL verification should be used for target group")
	cmd.Flags().Bool("sni", false, "Specifies SNI should be used for target group")

	return cmd
}

func loadbalancerTargetGroupUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	patchRequest := loadbalancer.PatchTargetGroupRequest{}
	patchRequest.Name, _ = cmd.Flags().GetString("name")
	patchRequest.Close = helper.GetBoolPtrFlagIfChanged(cmd, "close")
	patchRequest.Sticky = helper.GetBoolPtrFlagIfChanged(cmd, "sticky")
	patchRequest.CookieOpts, _ = cmd.Flags().GetString("cookie-opts")
	patchRequest.Source, _ = cmd.Flags().GetString("source")
	patchRequest.TimeoutsConnect, _ = cmd.Flags().GetInt("timeouts-connect")
	patchRequest.TimeoutsServer, _ = cmd.Flags().GetInt("timeouts-server")
	patchRequest.TimeoutsHTTPRequest, _ = cmd.Flags().GetInt("timeouts-http-request")
	patchRequest.TimeoutsCheck, _ = cmd.Flags().GetInt("timeouts-check")
	patchRequest.TimeoutsTunnel, _ = cmd.Flags().GetInt("timeouts-tunnel")
	patchRequest.CustomOptions, _ = cmd.Flags().GetString("custom-options")
	patchRequest.MonitorURL, _ = cmd.Flags().GetString("monitor-url")
	patchRequest.MonitorHost, _ = cmd.Flags().GetString("monitor-host")
	patchRequest.MonitorHTTPVersion, _ = cmd.Flags().GetString("monitor-http-version")
	patchRequest.MonitorExpect, _ = cmd.Flags().GetString("monitor-expect")
	patchRequest.MonitorExpectString, _ = cmd.Flags().GetString("monitor-expect-string")
	patchRequest.MonitorExpectStringRegex = helper.GetBoolPtrFlagIfChanged(cmd, "monitor-expect-string-regex")
	patchRequest.MonitorTCPMonitoring = helper.GetBoolPtrFlagIfChanged(cmd, "monitor-tcp-monitoring")
	patchRequest.CheckPort, _ = cmd.Flags().GetInt("check-port")
	patchRequest.SendProxy = helper.GetBoolPtrFlagIfChanged(cmd, "send-proxy")
	patchRequest.SendProxyV2 = helper.GetBoolPtrFlagIfChanged(cmd, "send-proxyv2")
	patchRequest.SSL = helper.GetBoolPtrFlagIfChanged(cmd, "ssl")
	patchRequest.SSLVerify = helper.GetBoolPtrFlagIfChanged(cmd, "ssl-verify")
	patchRequest.SNI = helper.GetBoolPtrFlagIfChanged(cmd, "sni")

	if cmd.Flags().Changed("balance") {
		balance, _ := cmd.Flags().GetString("balance")
		parsedBalance, err := loadbalancer.TargetGroupBalanceEnum.Parse(balance)
		if err != nil {
			return clierrors.NewErrInvalidFlagValue("balance", balance, err)
		}
		patchRequest.Balance = parsedBalance
	}

	if cmd.Flags().Changed("mode") {
		mode, _ := cmd.Flags().GetString("mode")
		parsedMode, err := loadbalancer.ModeEnum.Parse(mode)
		if err != nil {
			return clierrors.NewErrInvalidFlagValue("mode", mode, err)
		}
		patchRequest.Mode = parsedMode
	}

	if cmd.Flags().Changed("monitor-method") {
		monitorMethod, _ := cmd.Flags().GetString("monitor-method")
		parsedMonitorMethod, err := loadbalancer.TargetGroupMonitorMethodEnum.Parse(monitorMethod)
		if err != nil {
			return clierrors.NewErrInvalidFlagValue("monitor-method", monitorMethod, err)
		}
		patchRequest.MonitorMethod = parsedMonitorMethod
	}

	var groups []loadbalancer.TargetGroup
	for _, arg := range args {
		groupID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid target group ID [%s]", arg)
			continue
		}

		err = service.PatchTargetGroup(groupID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating target group [%d]: %s", groupID, err)
			continue
		}

		targetgroup, err := service.GetTargetGroup(groupID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated target group [%d]: %s", groupID, err)
			continue
		}

		groups = append(groups, targetgroup)
	}

	return output.CommandOutput(cmd, TargetGroupCollection(groups))
}

func loadbalancerTargetGroupDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <targetgroup: id>...",
		Short:   "Removes a target group",
		Long:    "This command removes one or more target groups",
		Example: "ans loadbalancer targetgroup delete 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing target group")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerTargetGroupDelete),
	}
}

func loadbalancerTargetGroupDelete(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		groupID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid target group ID [%s]", arg)
			continue
		}

		err = service.DeleteTargetGroup(groupID)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing target group [%d]: %s", groupID, err)
			continue
		}
	}

	return nil
}
