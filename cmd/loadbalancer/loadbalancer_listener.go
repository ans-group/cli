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
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func loadbalancerListenerRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listener",
		Short: "sub-commands relating to listeners",
	}

	// Child commands
	cmd.AddCommand(loadbalancerListenerListCmd(f))
	cmd.AddCommand(loadbalancerListenerShowCmd(f))
	cmd.AddCommand(loadbalancerListenerCreateCmd(f))
	cmd.AddCommand(loadbalancerListenerUpdateCmd(f))
	cmd.AddCommand(loadbalancerListenerDeleteCmd(f))

	// Child root commands
	cmd.AddCommand(loadbalancerListenerAccessIPRootCmd(f))
	cmd.AddCommand(loadbalancerListenerACLRootCmd(f))
	cmd.AddCommand(loadbalancerListenerCertificateRootCmd(f, fs))

	return cmd
}

func loadbalancerListenerListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists listeners",
		Long:    "This command lists listeners",
		Example: "ans loadbalancer listener list",
		RunE:    loadbalancerCobraRunEFunc(f, loadbalancerListenerList),
	}
}

func loadbalancerListenerList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	listeners, err := service.GetListeners(params)
	if err != nil {
		return fmt.Errorf("Error retrieving listeners: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerListenersProvider(listeners))
}

func loadbalancerListenerShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <listener: id>...",
		Short:   "Shows a listener",
		Long:    "This command shows one or more listeners",
		Example: "ans loadbalancer listener show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerShow),
	}
}

func loadbalancerListenerShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var listeners []loadbalancer.Listener
	for _, arg := range args {
		listenerID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid listener ID [%s]", arg)
			continue
		}

		listener, err := service.GetListener(listenerID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving listener [%d]: %s", listenerID, err)
			continue
		}

		listeners = append(listeners, listener)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerListenersProvider(listeners))
}

func loadbalancerListenerCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <listener: id>",
		Short:   "Creates a listener",
		Long:    "This command creates a listener",
		Example: "ans loadbalancer listener create --cluster 123 --default-target-group 456 --name \"test-listener\" --mode http",
		RunE:    loadbalancerCobraRunEFunc(f, loadbalancerListenerCreate),
	}

	cmd.Flags().String("name", "", "Name of listener")
	cmd.MarkFlagRequired("name")
	cmd.Flags().Int("cluster", 0, "ID of cluster")
	cmd.MarkFlagRequired("cluster")
	cmd.Flags().String("mode", "", "Specifies mode for listener")
	cmd.MarkFlagRequired("mode")
	cmd.Flags().Int("default-target-group", 0, "ID of default target group")
	cmd.MarkFlagRequired("default-target-group")
	cmd.Flags().Bool("hsts-enabled", false, "Specifies HSTS should be enabled for listener")
	cmd.Flags().Int("hsts-max-age", 0, "Specifies HSTS maximum age for listener")
	cmd.Flags().Bool("close", false, "Specifies close should be enabled for listener")
	cmd.Flags().Bool("redirect-https", false, "Specifies HTTPS redirection should be enabled")
	cmd.Flags().Bool("access-is-allow-list", false, "Specifies access IP behaviour should be allow rather than block")
	cmd.Flags().Bool("allow-tlsv1", false, "Specifies TLSv1 should be allowed")
	cmd.Flags().Bool("allow-tlsv11", false, "Specifies TLSv1.1 should be allowed")
	cmd.Flags().Bool("disable-tlsv12", false, "Specifies TLSv1.2 should be disabled")
	cmd.Flags().Bool("disable-http2", false, "Specifies HTTP2 should be disabled")
	cmd.Flags().String("custom-ciphers", "", "Specifies custom ciphers for listener")
	cmd.Flags().String("geoip-restriction", "", "Specifies restriction for GeoIP")
	cmd.Flags().StringSlice("geoip-continent", []string{""}, "Specifies continent for GeoIP. Can be repeated")
	cmd.Flags().StringSlice("geoip-country", []string{""}, "Specifies country for GeoIP. Can be repeated")
	cmd.Flags().Bool("geoip-european-union", false, "Specifies European Union for GeoIP. Can be repeated")

	return cmd
}

func loadbalancerListenerCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	createRequest := loadbalancer.CreateListenerRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.ClusterID, _ = cmd.Flags().GetInt("cluster")
	createRequest.HSTSEnabled, _ = cmd.Flags().GetBool("hsts-enabled")
	createRequest.HSTSMaxAge, _ = cmd.Flags().GetInt("hsts-max-age")
	createRequest.Close, _ = cmd.Flags().GetBool("close")
	createRequest.RedirectHTTPS, _ = cmd.Flags().GetBool("redirect-https")
	createRequest.DefaultTargetGroupID, _ = cmd.Flags().GetInt("default-target-group")
	createRequest.AccessIsAllowList, _ = cmd.Flags().GetBool("access-is-allow-list")
	createRequest.AllowTLSV1, _ = cmd.Flags().GetBool("allow-tlsv1")
	createRequest.AllowTLSV11, _ = cmd.Flags().GetBool("allow-tlsv11")
	createRequest.DisableTLSV12, _ = cmd.Flags().GetBool("disable-tlsv12")
	createRequest.DisableHTTP2, _ = cmd.Flags().GetBool("disable-http2")
	createRequest.CustomCiphers, _ = cmd.Flags().GetString("custom-ciphers")

	geoip := &loadbalancer.ListenerGeoIPRequest{}
	geoipChanged := false

	if cmd.Flags().Changed("geoip-restriction") {
		geoipChanged = true

		geoipRestrictionFlag, _ := cmd.Flags().GetString("geoip-restriction")
		geoipRestriction, err := loadbalancer.ParseListenerGeoIPRestriction(geoipRestrictionFlag)
		if err != nil {
			return clierrors.NewErrInvalidFlagValue("geoip-restriction", geoipRestrictionFlag, err)
		}

		geoip.Restriction = geoipRestriction
	}

	if cmd.Flags().Changed("geoip-continent") {
		geoipChanged = true
		geoip.Continents, _ = cmd.Flags().GetStringSlice("geoip-continent")
	}

	if cmd.Flags().Changed("geoip-country") {
		geoipChanged = true
		geoip.Countries, _ = cmd.Flags().GetStringSlice("geoip-country")
	}

	if cmd.Flags().Changed("geoip-european-union") {
		geoipChanged = true
		geoipEuropeanUnion, _ := cmd.Flags().GetBool("geoip-european-union")
		geoip.EuropeanUnion = &geoipEuropeanUnion
	}

	if geoipChanged {
		createRequest.GeoIP = geoip
	}

	mode, _ := cmd.Flags().GetString("mode")
	parsedMode, err := loadbalancer.ParseMode(mode)
	if err != nil {
		return clierrors.NewErrInvalidFlagValue("mode", mode, err)
	}
	createRequest.Mode = parsedMode

	listenerID, err := service.CreateListener(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating listener: %s", err)
	}

	listener, err := service.GetListener(listenerID)
	if err != nil {
		return fmt.Errorf("Error retrieving new listener: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerListenersProvider([]loadbalancer.Listener{listener}))
}

func loadbalancerListenerUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <listener: id>...",
		Short:   "Updates a listener",
		Long:    "This command updates one or more listeners",
		Example: "ans loadbalancer listener update 123 --name mylistener",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerUpdate),
	}

	cmd.Flags().String("name", "", "Name of listener")
	cmd.Flags().Bool("hsts-enabled", false, "Specifies HSTS should be enabled for listener")
	cmd.Flags().String("mode", "", "Specifies mode for listener")
	cmd.Flags().Int("hsts-max-age", 0, "Specifies HSTS maximum age for listener")
	cmd.Flags().Bool("close", false, "Specifies close should be enabled for listener")
	cmd.Flags().Bool("redirect-https", false, "Specifies HTTPS redirection should be enabled")
	cmd.Flags().Int("default-target-group", 0, "ID of default target group")
	cmd.Flags().Bool("access-is-allow-list", false, "Specifies access IP behaviour should be allow rather than block")
	cmd.Flags().Bool("allow-tlsv1", false, "Specifies TLSv1 should be allowed")
	cmd.Flags().Bool("allow-tlsv11", false, "Specifies TLSv1.1 should be allowed")
	cmd.Flags().Bool("disable-tlsv12", false, "Specifies TLSv1.2 should be disabled")
	cmd.Flags().Bool("disable-http2", false, "Specifies HTTP2 should be disabled")
	cmd.Flags().String("custom-ciphers", "", "Specifies custom ciphers for listener")
	cmd.Flags().Bool("geoip-disabled", false, "Specifies GeoIP should be disabled for listener")
	cmd.Flags().String("geoip-restriction", "", "Specifies restriction for GeoIP")
	cmd.Flags().StringSlice("geoip-continent", []string{""}, "Specifies continent for GeoIP. Can be repeated")
	cmd.Flags().StringSlice("geoip-country", []string{""}, "Specifies country for GeoIP. Can be repeated")
	cmd.Flags().Bool("geoip-european-union", false, "Specifies European Union for GeoIP. Can be repeated")

	return cmd
}

func loadbalancerListenerUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	patchRequest := loadbalancer.PatchListenerRequest{}
	patchRequest.Name, _ = cmd.Flags().GetString("name")
	patchRequest.HSTSEnabled = helper.GetBoolPtrFlagIfChanged(cmd, "hsts-enabled")
	patchRequest.HSTSMaxAge, _ = cmd.Flags().GetInt("hsts-max-age")
	patchRequest.Close = helper.GetBoolPtrFlagIfChanged(cmd, "close")
	patchRequest.RedirectHTTPS = helper.GetBoolPtrFlagIfChanged(cmd, "redirect-https")
	patchRequest.DefaultTargetGroupID, _ = cmd.Flags().GetInt("default-target-group")
	patchRequest.AccessIsAllowList = helper.GetBoolPtrFlagIfChanged(cmd, "access-is-allow-list")
	patchRequest.AllowTLSV1 = helper.GetBoolPtrFlagIfChanged(cmd, "allow-tlsv1")
	patchRequest.AllowTLSV11 = helper.GetBoolPtrFlagIfChanged(cmd, "allow-tlsv11")
	patchRequest.DisableTLSV12 = helper.GetBoolPtrFlagIfChanged(cmd, "disable-tlsv12")
	patchRequest.DisableHTTP2 = helper.GetBoolPtrFlagIfChanged(cmd, "disable-http2")
	patchRequest.CustomCiphers, _ = cmd.Flags().GetString("custom-ciphers")

	geoip := &loadbalancer.ListenerGeoIPRequest{}
	geoipChanged := false

	if cmd.Flags().Changed("geoip-restriction") {
		geoipChanged = true

		geoipRestrictionFlag, _ := cmd.Flags().GetString("geoip-restriction")
		geoipRestriction, err := loadbalancer.ParseListenerGeoIPRestriction(geoipRestrictionFlag)
		if err != nil {
			return clierrors.NewErrInvalidFlagValue("geoip-restriction", geoipRestrictionFlag, err)
		}

		geoip.Restriction = geoipRestriction
	}

	if cmd.Flags().Changed("geoip-continent") {
		geoipChanged = true
		geoip.Continents, _ = cmd.Flags().GetStringSlice("geoip-continent")
	}

	if cmd.Flags().Changed("geoip-country") {
		geoipChanged = true
		geoip.Countries, _ = cmd.Flags().GetStringSlice("geoip-country")
	}

	if cmd.Flags().Changed("geoip-european-union") {
		geoipChanged = true
		geoipEuropeanUnion, _ := cmd.Flags().GetBool("geoip-european-union")
		geoip.EuropeanUnion = &geoipEuropeanUnion
	}

	if geoipChanged {
		patchRequest.GeoIP = geoip
	}

	var listeners []loadbalancer.Listener
	for _, arg := range args {
		listenerID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid listener ID [%s]", arg)
			continue
		}

		if geoipDisabled, _ := cmd.Flags().GetBool("geoip-disabled"); geoipDisabled {
			err = service.DisableListenerGeoIP(listenerID)
			if err != nil {
				output.OutputWithErrorLevelf("Error disabling GeoIP for listener [%d]: %s", listenerID, err)
				continue
			}
		}

		err = service.PatchListener(listenerID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating listener [%d]: %s", listenerID, err)
			continue
		}

		listener, err := service.GetListener(listenerID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated listener [%d]: %s", listenerID, err)
			continue
		}

		listeners = append(listeners, listener)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerListenersProvider(listeners))
}

func loadbalancerListenerDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <listener: id>...",
		Short:   "Removes a listener",
		Long:    "This command removes one or more listeners",
		Example: "ans loadbalancer listener delete 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerDelete),
	}
}

func loadbalancerListenerDelete(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		listenerID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid listener ID [%s]", arg)
			continue
		}

		err = service.DeleteListener(listenerID)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing listener [%d]: %s", listenerID, err)
			continue
		}
	}

	return nil
}
