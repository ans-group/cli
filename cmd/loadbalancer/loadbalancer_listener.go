package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func loadbalancerListenerRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listener",
		Short: "sub-commands relating to listeners",
	}

	// Child commands
	cmd.AddCommand(loadbalancerListenerListCmd(f))
	cmd.AddCommand(loadbalancerListenerShowCmd(f))
	cmd.AddCommand(loadbalancerListenerUpdateCmd(f))

	return cmd
}

func loadbalancerListenerListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists listeners",
		Long:    "This command lists listeners",
		Example: "ukfast loadbalancer listener list",
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
		Example: "ukfast loadbalancer listener show 123",
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
			output.OutputWithErrorLevelf("Error retrieving listener [%s]: %s", arg, err)
			continue
		}

		listeners = append(listeners, listener)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerListenersProvider(listeners))
}

func loadbalancerListenerCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <listener: id>...",
		Short:   "Creates a listener",
		Long:    "This command creates a listener",
		Example: "ukfast loadbalancer listener create --name mylistener",
		RunE:    loadbalancerCobraRunEFunc(f, loadbalancerListenerCreate),
	}

	cmd.Flags().String("name", "", "Name of listener")

	return cmd
}

func loadbalancerListenerCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	createRequest := loadbalancer.CreateListenerRequest{}

	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}

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
		Example: "ukfast loadbalancer listener update 123 --name mylistener",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerUpdate),
	}

	cmd.Flags().String("name", "", "Name of listener")

	return cmd
}

func loadbalancerListenerUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	patchRequest := loadbalancer.PatchListenerRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var listeners []loadbalancer.Listener
	for _, arg := range args {
		listenerID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid listener ID [%s]", arg)
			continue
		}

		err = service.PatchListener(listenerID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating listener [%s]: %s", arg, err)
			continue
		}

		listener, err := service.GetListener(listenerID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated listener [%s]: %s", arg, err)
			continue
		}

		listeners = append(listeners, listener)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerListenersProvider(listeners))
}
