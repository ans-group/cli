package ecloud_v2

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

func Test_ecloudRouterList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRouters(gomock.Any()).Return([]ecloud.Router{}, nil).Times(1)

		ecloudRouterList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudRouterList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetRoutersError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRouters(gomock.Any()).Return([]ecloud.Router{}, errors.New("test error")).Times(1)

		err := ecloudRouterList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving routers: test error", err.Error())
	})
}

func Test_ecloudRouterShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudRouterShowCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudRouterShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing router", err.Error())
	})
}

func Test_ecloudRouterShow(t *testing.T) {
	t.Run("SingleRouter", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRouter("rtr-abcdef12").Return(ecloud.Router{}, nil).Times(1)

		ecloudRouterShow(service, &cobra.Command{}, []string{"rtr-abcdef12"})
	})

	t.Run("MultipleRouters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetRouter("rtr-abcdef12").Return(ecloud.Router{}, nil),
			service.EXPECT().GetRouter("rtr-abcdef23").Return(ecloud.Router{}, nil),
		)

		ecloudRouterShow(service, &cobra.Command{}, []string{"rtr-abcdef12", "rtr-abcdef23"})
	})

	t.Run("GetRouterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRouter("rtr-abcdef12").Return(ecloud.Router{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving router [rtr-abcdef12]: test error\n", func() {
			ecloudRouterShow(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})
}
