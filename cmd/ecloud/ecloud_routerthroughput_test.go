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

func Test_ecloudRouterThroughputList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRouterThroughputs(gomock.Any()).Return([]ecloud.RouterThroughput{}, nil).Times(1)

		ecloudRouterThroughputList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudRouterThroughputList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetRouterThroughputsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRouterThroughputs(gomock.Any()).Return([]ecloud.RouterThroughput{}, errors.New("test error")).Times(1)

		err := ecloudRouterThroughputList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving router throughputs: test error", err.Error())
	})
}

func Test_ecloudRouterThroughputShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudRouterThroughputShowCmd(nil).Args(nil, []string{"rtp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudRouterThroughputShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing router throughput", err.Error())
	})
}

func Test_ecloudRouterThroughputShow(t *testing.T) {
	t.Run("SingleRouterThroughput", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRouterThroughput("rtp-abcdef12").Return(ecloud.RouterThroughput{}, nil).Times(1)

		ecloudRouterThroughputShow(service, &cobra.Command{}, []string{"rtp-abcdef12"})
	})

	t.Run("MultipleRouterThroughputs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetRouterThroughput("rtp-abcdef12").Return(ecloud.RouterThroughput{}, nil),
			service.EXPECT().GetRouterThroughput("rtp-abcdef23").Return(ecloud.RouterThroughput{}, nil),
		)

		ecloudRouterThroughputShow(service, &cobra.Command{}, []string{"rtp-abcdef12", "rtp-abcdef23"})
	})

	t.Run("GetRouterThroughputError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRouterThroughput("rtp-abcdef12").Return(ecloud.RouterThroughput{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving router throughput [rtp-abcdef12]: test error\n", func() {
			ecloudRouterThroughputShow(service, &cobra.Command{}, []string{"rtp-abcdef12"})
		})
	})
}
