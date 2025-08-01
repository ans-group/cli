package ecloud

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/ptr"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
)

func withDryRun(message string, fn func() error, dryRun bool) error {
	if dryRun {
		fmt.Printf("Would %s\n", message)
		return nil
	}
	return fn()
}

func confirmVPCRecursiveDeletion(vpcID string) (bool, error) {
	fmt.Printf("WARNING: This will recursively delete ALL resources within VPC [%s]\n", vpcID)
	fmt.Printf("This includes instances, load balancers, networks, routers, volumes, and more.\n")
	fmt.Printf("This action cannot be undone.\n")
	fmt.Printf("\n")
	fmt.Printf("To confirm, please re-enter the VPC ID: %s\n", vpcID)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("ecloud: Failed to read input: %s", err)
	}

	input = strings.TrimSpace(input)
	if input != vpcID {
		return false, nil
	}

	return true, nil
}

func deleteVPCResourcesRecursively(service ecloud.ECloudService, vpcID string, dryRun bool) error {
	if dryRun {
		fmt.Printf("DRY RUN: Showing resources that would be deleted in VPC [%s]\n", vpcID)
	} else {
		fmt.Printf("Starting recursive deletion of resources in VPC [%s]\n", vpcID)
	}

	err := deleteAffinityRules(service, vpcID, dryRun)
	if err != nil {
		return fmt.Errorf("failed to delete affinity rules: %s", err)
	}

	err = deleteInstances(service, vpcID, dryRun)
	if err != nil {
		return fmt.Errorf("failed to delete instances: %s", err)
	}

	err = deleteLoadBalancers(service, vpcID, dryRun)
	if err != nil {
		return fmt.Errorf("failed to delete load balancers: %s", err)
	}

	err = deleteVPNResources(service, vpcID, dryRun)
	if err != nil {
		return fmt.Errorf("failed to delete VPN resources: %s", err)
	}

	err = deleteNetworkResources(service, vpcID, dryRun)
	if err != nil {
		return fmt.Errorf("failed to delete network resources: %s", err)
	}

	err = deleteRemainingFloatingIPs(service, vpcID, dryRun)
	if err != nil {
		return fmt.Errorf("failed to delete remaining floating IPs: %s", err)
	}

	err = deleteVolumeResources(service, vpcID, dryRun)
	if err != nil {
		return fmt.Errorf("failed to delete volume resources: %s", err)
	}

	err = deleteHostResources(service, vpcID, dryRun)
	if err != nil {
		return fmt.Errorf("failed to delete host resources: %s", err)
	}

	err = deletePrivateImages(service, vpcID, dryRun)
	if err != nil {
		return fmt.Errorf("failed to delete private images: %s", err)
	}

	err = deleteVPC(service, vpcID, dryRun)
	if err != nil {
		return fmt.Errorf("failed to delete VPC: %s", err)
	}

	if dryRun {
		fmt.Printf("DRY RUN: Completed listing resources that would be deleted in VPC [%s]\n", vpcID)
	} else {
		fmt.Printf("Completed recursive deletion of resources in VPC [%s]\n", vpcID)
	}
	return nil
}

func deleteInstances(service ecloud.ECloudService, vpcID string, dryRun bool) error {
	if dryRun {
		fmt.Printf("Would delete instances...\n")
	} else {
		fmt.Printf("Deleting instances...\n")
	}

	instances, err := service.GetVPCInstances(vpcID, connection.APIRequestParameters{})
	if err != nil {
		return fmt.Errorf("failed to get instances: %s", err)
	}

	for _, instance := range instances {
		err := withDryRun(fmt.Sprintf("delete instance [%s]", instance.ID), func() error {
			fmt.Printf("Deleting instance [%s]\n", instance.ID)
			err := service.DeleteInstance(instance.ID)
			if err != nil {
				output.OutputWithErrorLevelf("Error deleting instance [%s]: %s", instance.ID, err)
				return nil // Continue with other instances
			}

			err = helper.WaitForCommand(InstanceNotFoundWaitFunc(service, instance.ID))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for instance [%s] deletion: %s", instance.ID, err)
			}
			return nil
		}, dryRun)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteLoadBalancers(service ecloud.ECloudService, vpcID string, dryRun bool) error {
	if dryRun {
		fmt.Printf("Would delete load balancers...\n")
	} else {
		fmt.Printf("Deleting load balancers...\n")
	}

	params := connection.APIRequestParameters{
		Filtering: []connection.APIRequestFiltering{
			{
				Property: "vpc_id",
				Operator: connection.EQOperator,
				Value:    []string{vpcID},
			},
		},
	}

	loadBalancers, err := service.GetLoadBalancers(params)
	if err != nil {
		return fmt.Errorf("failed to get load balancers: %s", err)
	}

	for _, lb := range loadBalancers {
		err := withDryRun(fmt.Sprintf("delete load balancer [%s]", lb.ID), func() error {
			fmt.Printf("Deleting load balancer [%s]\n", lb.ID)
			taskID, err := service.DeleteLoadBalancer(lb.ID)
			if err != nil {
				output.OutputWithErrorLevelf("Error deleting load balancer [%s]: %s", lb.ID, err)
				return nil // Continue with other load balancers
			}

			err = helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for load balancer [%s] deletion: %s", lb.ID, err)
			}
			return nil
		}, dryRun)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteVPNResources(service ecloud.ECloudService, vpcID string, dryRun bool) error {
	if dryRun {
		fmt.Printf("Would delete VPN resources...\n")
	} else {
		fmt.Printf("Deleting VPN resources...\n")
	}

	params := connection.APIRequestParameters{
		Filtering: []connection.APIRequestFiltering{
			{
				Property: "vpc_id",
				Operator: connection.EQOperator,
				Value:    []string{vpcID},
			},
		},
	}

	vpnServices, err := service.GetVPNServices(params)
	if err != nil {
		return fmt.Errorf("failed to get VPN services: %s", err)
	}

	for _, vpnService := range vpnServices {
		// Delete VPN sessions first
		sessionParams := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "vpn_service_id",
					Operator: connection.EQOperator,
					Value:    []string{vpnService.ID},
				},
			},
		}

		vpnSessions, err := service.GetVPNSessions(sessionParams)
		if err != nil {
			output.OutputWithErrorLevelf("Error getting VPN sessions for service [%s]: %s", vpnService.ID, err)
		} else {
			for _, session := range vpnSessions {
				err := withDryRun(fmt.Sprintf("delete VPN session [%s]", session.ID), func() error {
					fmt.Printf("Deleting VPN session [%s]\n", session.ID)
					taskID, err := service.DeleteVPNSession(session.ID)
					if err != nil {
						output.OutputWithErrorLevelf("Error deleting VPN session [%s]: %s", session.ID, err)
						return nil // Continue with other sessions
					}

					err = helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
					if err != nil {
						output.OutputWithErrorLevelf("Error waiting for VPN session [%s] deletion: %s", session.ID, err)
					}
					return nil
				}, dryRun)
				if err != nil {
					return err
				}
			}
		}

		// Delete VPN endpoints
		endpointParams := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "vpn_service_id",
					Operator: connection.EQOperator,
					Value:    []string{vpnService.ID},
				},
			},
		}

		vpnEndpoints, err := service.GetVPNEndpoints(endpointParams)
		if err != nil {
			output.OutputWithErrorLevelf("Error getting VPN endpoints for service [%s]: %s", vpnService.ID, err)
		} else {
			for _, endpoint := range vpnEndpoints {
				err := withDryRun(fmt.Sprintf("delete VPN endpoint [%s]", endpoint.ID), func() error {
					fmt.Printf("Deleting VPN endpoint [%s]\n", endpoint.ID)
					taskID, err := service.DeleteVPNEndpoint(endpoint.ID)
					if err != nil {
						output.OutputWithErrorLevelf("Error deleting VPN endpoint [%s]: %s", endpoint.ID, err)
						return nil // Continue with other endpoints
					}

					err = helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
					if err != nil {
						output.OutputWithErrorLevelf("Error waiting for VPN endpoint [%s] deletion: %s", endpoint.ID, err)
					}
					return nil
				}, dryRun)
				if err != nil {
					return err
				}
			}
		}

		// Finally delete the VPN service
		err = withDryRun(fmt.Sprintf("delete VPN service [%s]", vpnService.ID), func() error {
			fmt.Printf("Deleting VPN service [%s]\n", vpnService.ID)
			taskID, err := service.DeleteVPNService(vpnService.ID)
			if err != nil {
				output.OutputWithErrorLevelf("Error deleting VPN service [%s]: %s", vpnService.ID, err)
				return nil // Continue with other services
			}

			err = helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for VPN service [%s] deletion: %s", vpnService.ID, err)
			}
			return nil
		}, dryRun)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteNetworkResources(service ecloud.ECloudService, vpcID string, dryRun bool) error {
	if dryRun {
		fmt.Printf("Would delete network resources...\n")
	} else {
		fmt.Printf("Deleting network resources...\n")
	}

	routerParams := connection.APIRequestParameters{
		Filtering: []connection.APIRequestFiltering{
			{
				Property: "vpc_id",
				Operator: connection.EQOperator,
				Value:    []string{vpcID},
			},
		},
	}

	routers, err := service.GetRouters(routerParams)
	if err != nil {
		return fmt.Errorf("failed to get routers: %s", err)
	}

	for _, router := range routers {
		// Get networks for this router
		networkParams := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "router_id",
					Operator: connection.EQOperator,
					Value:    []string{router.ID},
				},
			},
		}

		networks, err := service.GetNetworks(networkParams)
		if err != nil {
			output.OutputWithErrorLevelf("Error getting networks for router [%s]: %s", router.ID, err)
			continue
		}

		for _, network := range networks {
			// Delete cluster IP addresses
			ipParams := connection.APIRequestParameters{
				Filtering: []connection.APIRequestFiltering{
					{
						Property: "network_id",
						Operator: connection.EQOperator,
						Value:    []string{network.ID},
					},
					{
						Property: "type",
						Operator: connection.EQOperator,
						Value:    []string{"cluster"},
					},
				},
			}

			ipAddresses, err := service.GetIPAddresses(ipParams)
			if err != nil {
				output.OutputWithErrorLevelf("Error getting IP addresses for network [%s]: %s", network.ID, err)
			} else {
				for _, ip := range ipAddresses {
					if dryRun {
						fmt.Printf("Would delete IP address [%s]\n", ip.ID)
					} else {
						fmt.Printf("Deleting IP address [%s]\n", ip.ID)
						taskID, err := service.DeleteIPAddress(ip.ID)
						if err != nil {
							output.OutputWithErrorLevelf("Error deleting IP address [%s]: %s", ip.ID, err)
							continue
						}

						err = helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
						if err != nil {
							output.OutputWithErrorLevelf("Error waiting for IP address [%s] deletion: %s", ip.ID, err)
						}
					}
				}
			}

			// Delete NAT overload rules
			natParams := connection.APIRequestParameters{
				Filtering: []connection.APIRequestFiltering{
					{
						Property: "network_id",
						Operator: connection.EQOperator,
						Value:    []string{network.ID},
					},
				},
			}

			natRules, err := service.GetNATOverloadRules(natParams)
			if err != nil {
				output.OutputWithErrorLevelf("Error getting NAT overload rules for network [%s]: %s", network.ID, err)
			} else {
				for _, rule := range natRules {
					if dryRun {
						fmt.Printf("Would delete NAT overload rule [%s]\n", rule.ID)
					} else {
						fmt.Printf("Deleting NAT overload rule [%s]\n", rule.ID)
						taskID, err := service.DeleteNATOverloadRule(rule.ID)
						if err != nil {
							output.OutputWithErrorLevelf("Error deleting NAT overload rule [%s]: %s", rule.ID, err)
							continue
						}

						err = helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
						if err != nil {
							output.OutputWithErrorLevelf("Error waiting for NAT overload rule [%s] deletion: %s", rule.ID, err)
						}
					}
				}
			}

			// Delete the network
			if dryRun {
				fmt.Printf("Would delete network [%s]\n", network.ID)
			} else {
				fmt.Printf("Deleting network [%s]\n", network.ID)
				err = service.DeleteNetwork(network.ID)
				if err != nil {
					output.OutputWithErrorLevelf("Error deleting network [%s]: %s", network.ID, err)
				} else {
					err = helper.WaitForCommand(NetworkNotFoundWaitFunc(service, network.ID))
					if err != nil {
						output.OutputWithErrorLevelf("Error waiting for network [%s] deletion: %s", network.ID, err)
					}
				}
			}
		}

		// Unassign floating IPs from router
		fipParams := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "resource_id",
					Operator: connection.EQOperator,
					Value:    []string{router.ID},
				},
			},
		}

		floatingIPs, err := service.GetFloatingIPs(fipParams)
		if err != nil {
			output.OutputWithErrorLevelf("Error getting floating IPs for router [%s]: %s", router.ID, err)
		} else {
			for _, fip := range floatingIPs {
				if dryRun {
					fmt.Printf("Would unassign floating IP [%s] from router [%s]\n", fip.ID, router.ID)
				} else {
					fmt.Printf("Unassigning floating IP [%s] from router [%s]\n", fip.ID, router.ID)
					taskID, err := service.UnassignFloatingIP(fip.ID)
					if err != nil {
						output.OutputWithErrorLevelf("Error unassigning floating IP [%s]: %s", fip.ID, err)
						continue
					}

					err = helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
					if err != nil {
						output.OutputWithErrorLevelf("Error waiting for floating IP [%s] unassignment: %s", fip.ID, err)
					}
				}
			}
		}

		// Delete the router
		if dryRun {
			fmt.Printf("Would delete router [%s]\n", router.ID)
		} else {
			fmt.Printf("Deleting router [%s]\n", router.ID)
			err = service.DeleteRouter(router.ID)
			if err != nil {
				output.OutputWithErrorLevelf("Error deleting router [%s]: %s", router.ID, err)
			} else {
				err = helper.WaitForCommand(RouterNotFoundWaitFunc(service, router.ID))
				if err != nil {
					output.OutputWithErrorLevelf("Error waiting for router [%s] deletion: %s", router.ID, err)
				}
			}
		}
	}

	return nil
}

func deleteRemainingFloatingIPs(service ecloud.ECloudService, vpcID string, dryRun bool) error {
	if dryRun {
		fmt.Printf("Would delete remaining floating IPs...\n")
	} else {
		fmt.Printf("Deleting remaining floating IPs...\n")
	}

	params := connection.APIRequestParameters{
		Filtering: []connection.APIRequestFiltering{
			{
				Property: "vpc_id",
				Operator: connection.EQOperator,
				Value:    []string{vpcID},
			},
		},
	}

	floatingIPs, err := service.GetFloatingIPs(params)
	if err != nil {
		return fmt.Errorf("failed to get floating IPs: %s", err)
	}

	for _, fip := range floatingIPs {
		err := withDryRun(fmt.Sprintf("delete floating IP [%s]", fip.ID), func() error {
			fmt.Printf("Deleting floating IP [%s]\n", fip.ID)
			taskID, err := service.DeleteFloatingIP(fip.ID)
			if err != nil {
				output.OutputWithErrorLevelf("Error deleting floating IP [%s]: %s", fip.ID, err)
				return nil // Continue with other floating IPs
			}

			err = helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for floating IP [%s] deletion: %s", fip.ID, err)
			}
			return nil
		}, dryRun)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteVolumeResources(service ecloud.ECloudService, vpcID string, dryRun bool) error {
	if dryRun {
		fmt.Printf("Would delete volume resources...\n")
	} else {
		fmt.Printf("Deleting volume resources...\n")
	}

	// First, remove volumes from volume groups (like bash script does)
	volumes, err := service.GetVPCVolumes(vpcID, connection.APIRequestParameters{})
	if err != nil {
		return fmt.Errorf("failed to get volumes: %s", err)
	}

	for _, volume := range volumes {
		if volume.VolumeGroupID != "" {
			err := withDryRun(fmt.Sprintf("remove volume [%s] from volume group", volume.ID), func() error {
				fmt.Printf("Removing volume [%s] from volume group\n", volume.ID)
				// Update volume to remove it from volume group
				patchRequest := ecloud.PatchVolumeRequest{
					VolumeGroupID: ptr.String(""), // Empty string to remove from group
				}
				task, err := service.PatchVolume(volume.ID, patchRequest)
				if err != nil {
					output.OutputWithErrorLevelf("Error removing volume [%s] from volume group: %s", volume.ID, err)
					return nil // Continue with other volumes
				}

				err = helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
				if err != nil {
					output.OutputWithErrorLevelf("Error waiting for volume [%s] group removal: %s", volume.ID, err)
				}
				return nil
			}, dryRun)
			if err != nil {
				return err
			}
		}
	}

	// Delete volumes
	for _, volume := range volumes {
		err := withDryRun(fmt.Sprintf("delete volume [%s]", volume.ID), func() error {
			fmt.Printf("Deleting volume [%s]\n", volume.ID)
			taskID, err := service.DeleteVolume(volume.ID)
			if err != nil {
				output.OutputWithErrorLevelf("Error deleting volume [%s]: %s", volume.ID, err)
				return nil // Continue with other volumes
			}

			err = helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for volume [%s] deletion: %s", volume.ID, err)
			}
			return nil
		}, dryRun)
		if err != nil {
			return err
		}
	}

	// Delete volume groups
	params := connection.APIRequestParameters{
		Filtering: []connection.APIRequestFiltering{
			{
				Property: "vpc_id",
				Operator: connection.EQOperator,
				Value:    []string{vpcID},
			},
		},
	}

	volumeGroups, err := service.GetVolumeGroups(params)
	if err != nil {
		return fmt.Errorf("failed to get volume groups: %s", err)
	}

	for _, vg := range volumeGroups {
		err := withDryRun(fmt.Sprintf("delete volume group [%s]", vg.ID), func() error {
			fmt.Printf("Deleting volume group [%s]\n", vg.ID)
			taskID, err := service.DeleteVolumeGroup(vg.ID)
			if err != nil {
				output.OutputWithErrorLevelf("Error deleting volume group [%s]: %s", vg.ID, err)
				return nil // Continue with other volume groups
			}

			err = helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for volume group [%s] deletion: %s", vg.ID, err)
			}
			return nil
		}, dryRun)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteHostResources(service ecloud.ECloudService, vpcID string, dryRun bool) error {
	if dryRun {
		fmt.Printf("Would delete host resources...\n")
	} else {
		fmt.Printf("Deleting host resources...\n")
	}

	params := connection.APIRequestParameters{
		Filtering: []connection.APIRequestFiltering{
			{
				Property: "vpc_id",
				Operator: connection.EQOperator,
				Value:    []string{vpcID},
			},
		},
	}

	hostGroups, err := service.GetHostGroups(params)
	if err != nil {
		return fmt.Errorf("failed to get host groups: %s", err)
	}

	for _, hostGroup := range hostGroups {
		// Delete hosts in the group first
		hostParams := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "host_group_id",
					Operator: connection.EQOperator,
					Value:    []string{hostGroup.ID},
				},
			},
		}

		hosts, err := service.GetHosts(hostParams)
		if err != nil {
			output.OutputWithErrorLevelf("Error getting hosts for host group [%s]: %s", hostGroup.ID, err)
		} else {
			for _, host := range hosts {
				if dryRun {
					fmt.Printf("Would delete host [%s]\n", host.ID)
				} else {
					fmt.Printf("Deleting host [%s]\n", host.ID)
					taskID, err := service.DeleteHost(host.ID)
					if err != nil {
						output.OutputWithErrorLevelf("Error deleting host [%s]: %s", host.ID, err)
						continue
					}

					err = helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
					if err != nil {
						output.OutputWithErrorLevelf("Error waiting for host [%s] deletion: %s", host.ID, err)
					}
				}
			}
		}

		// Delete the host group
		if dryRun {
			fmt.Printf("Would delete host group [%s]\n", hostGroup.ID)
		} else {
			fmt.Printf("Deleting host group [%s]\n", hostGroup.ID)
			taskID, err := service.DeleteHostGroup(hostGroup.ID)
			if err != nil {
				output.OutputWithErrorLevelf("Error deleting host group [%s]: %s", hostGroup.ID, err)
				continue
			}

			err = helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for host group [%s] deletion: %s", hostGroup.ID, err)
			}
		}
	}

	return nil
}

func deletePrivateImages(service ecloud.ECloudService, vpcID string, dryRun bool) error {
	if dryRun {
		fmt.Printf("Would delete private images...\n")
	} else {
		fmt.Printf("Deleting private images...\n")
	}

	params := connection.APIRequestParameters{
		Filtering: []connection.APIRequestFiltering{
			{
				Property: "vpc_id",
				Operator: connection.EQOperator,
				Value:    []string{vpcID},
			},
		},
	}

	images, err := service.GetImages(params)
	if err != nil {
		return fmt.Errorf("failed to get images: %s", err)
	}

	for _, image := range images {
		err := withDryRun(fmt.Sprintf("delete private image [%s]", image.ID), func() error {
			fmt.Printf("Deleting private image [%s]\n", image.ID)
			taskID, err := service.DeleteImage(image.ID)
			if err != nil {
				output.OutputWithErrorLevelf("Error deleting image [%s]: %s", image.ID, err)
				return nil // Continue with other images
			}

			err = helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for image [%s] deletion: %s", image.ID, err)
			}
			return nil
		}, dryRun)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteVPC(service ecloud.ECloudService, vpcID string, dryRun bool) error {
	if dryRun {
		fmt.Printf("Would delete VPC [%s]\n", vpcID)
		return nil
	}

	fmt.Printf("Deleting VPC [%s]...\n", vpcID)

	// Check if VPC still exists before trying to delete
	_, err := service.GetVPC(vpcID)
	if err != nil {
		// If we can't get the VPC, it might already be deleted
		output.OutputWithErrorLevelf("VPC [%s] may already be deleted or inaccessible: %s", vpcID, err)
		return nil
	}

	err = service.DeleteVPC(vpcID)
	if err != nil {
		return fmt.Errorf("error deleting VPC [%s]: %s", vpcID, err)
	}

	err = helper.WaitForCommand(VPCNotFoundWaitFunc(service, vpcID))
	if err != nil {
		return fmt.Errorf("error waiting for VPC [%s] deletion: %s", vpcID, err)
	}

	fmt.Printf("Successfully deleted VPC [%s]\n", vpcID)
	return nil
}

func deleteAffinityRules(service ecloud.ECloudService, vpcID string, dryRun bool) error {
	if dryRun {
		fmt.Printf("Would delete affinity rules...\n")
	} else {
		fmt.Printf("Deleting affinity rules...\n")
	}

	params := connection.APIRequestParameters{
		Filtering: []connection.APIRequestFiltering{
			{
				Property: "vpc_id",
				Operator: connection.EQOperator,
				Value:    []string{vpcID},
			},
		},
	}

	affinityRules, err := service.GetAffinityRules(params)
	if err != nil {
		return fmt.Errorf("failed to get affinity rules: %s", err)
	}

	for _, rule := range affinityRules {
		err := withDryRun(fmt.Sprintf("delete affinity rule [%s]", rule.ID), func() error {
			fmt.Printf("Deleting affinity rule [%s]\n", rule.ID)
			taskID, err := service.DeleteAffinityRule(rule.ID)
			if err != nil {
				output.OutputWithErrorLevelf("Error deleting affinity rule [%s]: %s", rule.ID, err)
				return nil // Continue with other affinity rules
			}

			err = helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for affinity rule [%s] deletion: %s", rule.ID, err)
			}
			return nil
		}, dryRun)
		if err != nil {
			return err
		}
	}

	return nil
}
