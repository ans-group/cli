package loadbalancer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/spf13/cobra"
)

const clusterIDsOutputKey string = "loadbalancer_cluster_ids"

type clusterIDList []int

func (c *clusterIDList) UnmarshalJSON(data []byte) error {
	var v []any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	for _, item := range v {
		switch v := item.(type) {
		case int:
			*c = append(*c, v)
		case float64:
			*c = append(*c, int(v))
		case string:
			if intValue, err := strconv.Atoi(v); err == nil {
				*c = append(*c, intValue)
			} else {
				return err
			}
		default:
			return fmt.Errorf("unsupported type: %T", item)
		}
	}

	return nil
}

func loadbalancerTerraformCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "terraform",
		Short: "Terraform wrapper for handling deployments",
		Long: "Wraps Terraform to automatically deploy your changes after a successful apply. You must " +
			fmt.Sprintf("ensure you have configured a Terraform output called '%s' which ", clusterIDsOutputKey) +
			"contains the IDs of your loadbalancer clusters. After applying your Terraform configuration, " +
			"we'll deploy your staged changes to the loadbalancer.",
		Example:            "ans loadbalancer terraform apply",
		DisableFlagParsing: true,
		RunE:               loadbalancerCobraRunEFunc(f, loadbalancerTerraform),
	}

	return cmd
}

func loadbalancerTerraform(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	// Because we disable flag parsing, we'll need to handle --help ourselves
	if slices.Contains(args, "-h") || slices.Contains(args, "--help") || len(args) == 0 {
		return cmd.Help()
	}

	binPath, err := findTerraformBinary()
	if err != nil {
		return err
	}

	if err := runTerraform(binPath, args); err != nil {
		return err
	}

	// If we're not applying the configuration then exit now
	if !slices.Contains(args, "apply") {
		return nil
	}

	fmt.Printf("\nTerraform run complete, deploying the configuration to the loadbalancer...\n")

	clusterIDs, err := getClusterIDs(binPath)
	if err != nil {
		return err
	}

	fmt.Printf("Deploying cluster(s): %s\n\n",
		strings.Trim(strings.Join(strings.Fields(fmt.Sprint(clusterIDs)), ", "), "[]"))

	haveErrors := false
	for _, clusterID := range clusterIDs {
		fmt.Printf("Deploying %d... ", clusterID)
		err = service.DeployCluster(clusterID)
		if err != nil {
			haveErrors = true
			fmt.Printf("failed to deploy cluster: %d: %s\n", clusterID, err)
		} else {
			fmt.Printf("ok\n\n")
		}
	}

	if haveErrors {
		return fmt.Errorf("\nans: some clusters failed to deploy, please see above output")
	}

	return nil
}

func findTerraformBinary() (string, error) {
	path, err := exec.LookPath("terraform")
	if err != nil {
		path, err = exec.LookPath("tofu")
		if err != nil {
			return "", fmt.Errorf("ans: no terraform or tofu binaries found in path")
		}
	}
	return path, nil
}

func runTerraform(binPath string, args []string) error {
	tfCmd := exec.Command(binPath, args...)
	tfCmd.Stdout = os.Stdout
	tfCmd.Stderr = os.Stderr
	tfCmd.Stdin = os.Stdin

	if err := tfCmd.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return fmt.Errorf("ans: %s exited with %d, aborting deployment", binPath, exitError.ExitCode())
		}
		return fmt.Errorf("ans: error running %s: %s", binPath, err)
	}

	return nil
}

func getClusterIDs(binPath string) ([]int, error) {
	// Get the cluster ID from Terraform
	output, err := exec.Command(binPath, "output", "-json", clusterIDsOutputKey).CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", string(output))
		return nil, fmt.Errorf("ans: deployment failed: failed to get %s output from terraform, cannot deploy: %s",
			clusterIDsOutputKey, err)
	}

	var clusterIDs clusterIDList
	err = json.Unmarshal(output, &clusterIDs)
	if err != nil {
		return nil, fmt.Errorf("ans: deployment failed: failed to unmarshal Terraform output from key '%s': %s",
			clusterIDsOutputKey, err)
	}

	if len(clusterIDs) == 0 {
		return nil, fmt.Errorf("ans: deployment failed: no cluster IDs found in Terraform output key '%s'",
			clusterIDsOutputKey)
	}

	return clusterIDs, nil
}
