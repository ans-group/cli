//go:build windows

package ecloud

import (
	"fmt"

	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

// ecloudInstanceSSHWithAuth is not supported on Windows
func ecloudInstanceSSHWithAuth(service ecloud.ECloudService, cmd *cobra.Command, instanceID, ipAddress string) error {
	return fmt.Errorf("ssh: --auto flag is not supported on Windows")
}
