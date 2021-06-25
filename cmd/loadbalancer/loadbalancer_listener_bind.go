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

func loadbalancerListenerBindRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind",
		Short: "sub-commands relating to binds",
	}

	// Child commands
	cmd.AddCommand(loadbalancerListenerBindListCmd(f))
	cmd.AddCommand(loadbalancerListenerBindShowCmd(f))
	cmd.AddCommand(loadbalancerListenerBindCreateCmd(f))
	cmd.AddCommand(loadbalancerListenerBindUpdateCmd(f))
	cmd.AddCommand(loadbalancerListenerBindDeleteCmd(f))

	return cmd
}

func loadbalancerListenerBindListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <listener: id>",
		Short:   "Lists binds",
		Long:    "This command lists binds",
		Example: "ukfast loadbalancer listener bind list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerBindList),
	}
}

func loadbalancerListenerBindList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	listenerID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid listener ID")
	}

	binds, err := service.GetListenerBinds(listenerID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving binds: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerBindsProvider(binds))
}

func loadbalancerListenerBindShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <listener: id> <bind: id>...",
		Short:   "Shows a bind",
		Long:    "This command shows one or more binds",
		Example: "ukfast loadbalancer listener bind show 123 345",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}
			if len(args) < 2 {
				return errors.New("Missing bind")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerBindShow),
	}
}

func loadbalancerListenerBindShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	listenerID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid listener ID")
	}

	var binds []loadbalancer.Bind
	for _, arg := range args[1:] {

		bindID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid bind ID [%s]", arg)
			continue
		}

		bind, err := service.GetListenerBind(listenerID, bindID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving bind [%d]: %s", bindID, err)
			continue
		}

		binds = append(binds, bind)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerBindsProvider(binds))
}

func loadbalancerListenerBindCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <listener: id> <bind: id>...",
		Short:   "Creates a bind",
		Long:    "This command creates a bind",
		Example: "ukfast loadbalancer listener bind create 123 --vip 456 --port 443",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerBindCreate),
	}

	cmd.Flags().Int("vip", 0, "ID of VIP")
	cmd.MarkFlagRequired("vip")
	cmd.Flags().Int("port", 0, "Port number for bind")
	cmd.MarkFlagRequired("port")

	return cmd
}

func loadbalancerListenerBindCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	listenerID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid listener ID")
	}

	createRequest := loadbalancer.CreateBindRequest{}
	createRequest.VIPID, _ = cmd.Flags().GetInt("vip")
	createRequest.Port, _ = cmd.Flags().GetInt("port")

	bindID, err := service.CreateListenerBind(listenerID, createRequest)
	if err != nil {
		return fmt.Errorf("Error creating bind: %s", err)
	}

	bind, err := service.GetListenerBind(listenerID, bindID)
	if err != nil {
		return fmt.Errorf("Error retrieving new bind: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerBindsProvider([]loadbalancer.Bind{bind}))
}

func loadbalancerListenerBindUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <listener: id> <bind: id>...",
		Short:   "Updates a bind",
		Long:    "This command updates one or more binds",
		Example: "ukfast loadbalancer listener bind update 123 --port 443",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}
			if len(args) < 2 {
				return errors.New("Missing bind")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerBindUpdate),
	}

	cmd.Flags().Int("vip", 0, "ID of VIP")
	cmd.Flags().Int("port", 0, "Port number for bind")

	return cmd
}

func loadbalancerListenerBindUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	listenerID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid listener ID")
	}

	patchRequest := loadbalancer.PatchBindRequest{}
	patchRequest.VIPID, _ = cmd.Flags().GetInt("vip")
	patchRequest.Port, _ = cmd.Flags().GetInt("port")

	var binds []loadbalancer.Bind
	for _, arg := range args[1:] {
		bindID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid bind ID [%s]", arg)
			continue
		}

		err = service.PatchListenerBind(listenerID, bindID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating bind [%d]: %s", bindID, err)
			continue
		}

		bind, err := service.GetListenerBind(listenerID, bindID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated bind [%d]: %s", bindID, err)
			continue
		}

		binds = append(binds, bind)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerBindsProvider(binds))
}

func loadbalancerListenerBindDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <listener: id> <bind: id>...",
		Short:   "Removes a bind",
		Long:    "This command removes one or more binds",
		Example: "ukfast loadbalancer listener bind delete 123 456",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}
			if len(args) < 2 {
				return errors.New("Missing bind")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerBindDelete),
	}
}

func loadbalancerListenerBindDelete(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	listenerID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid listener ID")
	}

	for _, arg := range args[1:] {
		bindID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid bind ID [%s]", arg)
			continue
		}

		err = service.DeleteListenerBind(listenerID, bindID)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing bind [%d]: %s", bindID, err)
			continue
		}
	}

	return nil
}
