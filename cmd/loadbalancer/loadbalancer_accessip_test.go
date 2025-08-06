package loadbalancer

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_loadbalancerAccessIPShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerAccessIPShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerAccessIPShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing access IP", err.Error())
	})
}

func Test_loadbalancerAccessIPShow(t *testing.T) {
	t.Run("SingleAccessIP", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetAccessIP(123).Return(loadbalancer.AccessIP{}, nil).Times(1)

		loadbalancerAccessIPShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleAccessIPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetAccessIP(123).Return(loadbalancer.AccessIP{}, nil),
			service.EXPECT().GetAccessIP(456).Return(loadbalancer.AccessIP{}, nil),
		)

		loadbalancerAccessIPShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetAccessIPID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid access IP ID [abc]\n", func() {
			loadbalancerAccessIPShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetAccessIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetAccessIP(123).Return(loadbalancer.AccessIP{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving access IP [123]: test error\n", func() {
			loadbalancerAccessIPShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_loadbalancerAccessIPUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerAccessIPUpdateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerAccessIPUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing access IP", err.Error())
	})
}

func Test_loadbalancerAccessIPUpdate(t *testing.T) {
	t.Run("SingleGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		cmd := loadbalancerAccessIPUpdateCmd(nil)
		cmd.ParseFlags([]string{"--ip=1.2.3.4"})

		req := loadbalancer.PatchAccessIPRequest{
			IP: "1.2.3.4",
		}

		gomock.InOrder(
			service.EXPECT().PatchAccessIP(123, req).Return(nil),
			service.EXPECT().GetAccessIP(123).Return(loadbalancer.AccessIP{}, nil),
		)

		loadbalancerAccessIPUpdate(service, cmd, []string{"123"})
	})

	t.Run("MultipleAccessIPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchAccessIP(123, gomock.Any()).Return(nil),
			service.EXPECT().GetAccessIP(123).Return(loadbalancer.AccessIP{}, nil),
			service.EXPECT().PatchAccessIP(456, gomock.Any()).Return(nil),
			service.EXPECT().GetAccessIP(456).Return(loadbalancer.AccessIP{}, nil),
		)

		loadbalancerAccessIPUpdate(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("PatchAccessIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().PatchAccessIP(123, gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating access IP [123]: test error\n", func() {
			loadbalancerAccessIPUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})

	t.Run("GetAccessIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchAccessIP(123, gomock.Any()).Return(nil),
			service.EXPECT().GetAccessIP(123).Return(loadbalancer.AccessIP{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated access IP [123]: test error\n", func() {
			loadbalancerAccessIPUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_loadbalancerAccessIPDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerAccessIPDeleteCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerAccessIPDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing access IP", err.Error())
	})
}

func Test_loadbalancerAccessIPDelete(t *testing.T) {
	t.Run("SingleGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteAccessIP(123).Return(nil).Times(1)

		loadbalancerAccessIPDelete(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleAccessIPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteAccessIP(123).Return(nil),
			service.EXPECT().DeleteAccessIP(456).Return(nil),
		)

		loadbalancerAccessIPDelete(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("DeleteAccessIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteAccessIP(123).Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing access IP [123]: test error\n", func() {
			loadbalancerAccessIPDelete(service, &cobra.Command{}, []string{"123"})
		})
	})
}
