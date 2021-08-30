package ecloud

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudVPNSessionCredentialListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNSessionCredentialListCmd(nil).Args(nil, []string{"vpns-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNSessionCredentialListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VPN session", err.Error())
	})
}

func Test_ecloudVPNSessionCredentialList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNSessionCredentials("vpns-abcdef12", gomock.Any()).Return([]ecloud.Credential{}, nil).Times(1)

		ecloudVPNSessionCredentialList(service, &cobra.Command{}, []string{"vpns-abcdef12"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVPNSessionCredentialList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetInstancesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNSessionCredentials("vpns-abcdef12", gomock.Any()).Return([]ecloud.Credential{}, errors.New("test error")).Times(1)

		err := ecloudVPNSessionCredentialList(service, &cobra.Command{}, []string{"vpns-abcdef12"})

		assert.Equal(t, "Error retrieving VPN session credentials: test error", err.Error())
	})
}
