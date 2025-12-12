//go:build !windows

package ecloud

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

// ecloudInstanceSSHWithAuth establishes SSH connection using sshpass for automatic password authentication
func ecloudInstanceSSHWithAuth(service ecloud.ECloudService, cmd *cobra.Command, instanceID, ipAddress string) error {
	if _, err := exec.LookPath("sshpass"); err != nil {
		return fmt.Errorf("ssh: sshpass is not installed, please install it to use the --auto option")
	}

	credential, err := selectCredential(service, cmd, instanceID)
	if err != nil {
		return err
	}

	user, _ := cmd.Flags().GetString("user")
	port, _ := cmd.Flags().GetInt("port")
	sshArgs, _ := cmd.Flags().GetString("args")

	sshUser := user
	if credential.Username != "" {
		sshUser = credential.Username
	}

	// Using -e flag to read password from SSHPASS environment variable (more secure than -p)
	sshCmd := exec.Command("sshpass", "-e", "ssh",
		fmt.Sprintf("%s@%s", sshUser, ipAddress),
		"-p", strconv.Itoa(port),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
	)

	if sshArgs != "" {
		sshCmd.Args = append(sshCmd.Args, sshArgs)
	}

	sshCmd.Env = append(os.Environ(), fmt.Sprintf("SSHPASS=%s", credential.Password))

	sshCmd.Stdout = os.Stdout
	sshCmd.Stdin = os.Stdin
	sshCmd.Stderr = os.Stderr

	if err := sshCmd.Start(); err != nil {
		return fmt.Errorf("ssh: failed to start ssh command: %w", err)
	}

	return sshCmd.Wait()
}

func selectCredential(service ecloud.ECloudService, cmd *cobra.Command, instanceID string) (*ecloud.Credential, error) {
	credentialName, _ := cmd.Flags().GetString("credential-name")
	user, _ := cmd.Flags().GetString("user")

	params := connection.APIRequestParameters{}
	if credentialName != "" {
		params.WithFilter(connection.APIRequestFiltering{
			Property: "name",
			Operator: connection.EQOperator,
			Value:    []string{credentialName},
		})
	}

	credentials, err := service.GetInstanceCredentials(instanceID, params)
	if err != nil {
		return nil, fmt.Errorf("ssh: failed to retrieve credentials: %w", err)
	}

	if len(credentials) == 0 {
		if credentialName != "" {
			return nil, fmt.Errorf("ssh: credential '%s' not found for instance %s", credentialName, instanceID)
		}
		return nil, fmt.Errorf("ssh: no credentials found for instance %s", instanceID)
	}

	// Selection logic:
	// 1. If credential-name flag provided, already filtered - use first result
	// 2. Otherwise, match username from --user flag
	// 3. Fallback to first available credential

	if credentialName != "" {
		return &credentials[0], nil
	}

	for i := range credentials {
		if credentials[i].Username == user {
			return &credentials[i], nil
		}
	}

	return &credentials[0], nil
}
