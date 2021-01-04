package ecloud

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudFloatingIPList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFloatingIPs(gomock.Any()).Return([]ecloud.FloatingIP{}, nil).Times(1)

		ecloudFloatingIPList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudFloatingIPList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetFloatingIPsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFloatingIPs(gomock.Any()).Return([]ecloud.FloatingIP{}, errors.New("test error")).Times(1)

		err := ecloudFloatingIPList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving floating IPs: test error", err.Error())
	})
}

func Test_ecloudFloatingIPShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFloatingIPShowCmd(nil).Args(nil, []string{"fip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFloatingIPShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing floating IP", err.Error())
	})
}

func Test_ecloudFloatingIPShow(t *testing.T) {
	t.Run("SingleFloatingIP", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil).Times(1)

		ecloudFloatingIPShow(service, &cobra.Command{}, []string{"fip-abcdef12"})
	})

	t.Run("MultipleFloatingIPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil),
			service.EXPECT().GetFloatingIP("fip-abcdef23").Return(ecloud.FloatingIP{}, nil),
		)

		ecloudFloatingIPShow(service, &cobra.Command{}, []string{"fip-abcdef12", "fip-abcdef23"})
	})

	t.Run("GetFloatingIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving floating IP [fip-abcdef12]: test error\n", func() {
			ecloudFloatingIPShow(service, &cobra.Command{}, []string{"fip-abcdef12"})
		})
	})
}
