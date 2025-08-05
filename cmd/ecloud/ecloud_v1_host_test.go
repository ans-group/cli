package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudV1HostList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetV1Hosts(gomock.Any()).Return([]ecloud.V1Host{}, nil).Times(1)

		ecloudV1HostList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudV1HostList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetV1HostsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetV1Hosts(gomock.Any()).Return([]ecloud.V1Host{}, errors.New("test error")).Times(1)

		err := ecloudV1HostList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving hosts: test error", err.Error())
	})
}

func Test_ecloudV1HostShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudV1HostShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudV1HostShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing host", err.Error())
	})
}

func Test_ecloudV1HostShow(t *testing.T) {
	t.Run("SingleHost", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetV1Host(123).Return(ecloud.V1Host{}, nil).Times(1)

		ecloudV1HostShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleHosts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetV1Host(123).Return(ecloud.V1Host{}, nil),
			service.EXPECT().GetV1Host(456).Return(ecloud.V1Host{}, nil),
		)

		ecloudV1HostShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetV1HostID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid host ID [abc]\n", func() {
			ecloudV1HostShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetV1HostError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetV1Host(123).Return(ecloud.V1Host{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving host [123]: test error\n", func() {
			ecloudV1HostShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
