package loadbalancer

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_loadbalancerListenerCertificateListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerListenerCertificateListCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerCertificateListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing listener", err.Error())
	})
}

func Test_loadbalancerListenerCertificateList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetListenerCertificates(123, gomock.Any()).Return([]loadbalancer.Certificate{}, nil).Times(1)

		loadbalancerListenerCertificateList(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerCertificateListCmd(nil)
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadbalancerListenerCertificateList(service, cmd, []string{"123"})

		assert.NotNil(t, err)
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetListenerCertificatesError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetListenerCertificates(123, gomock.Any()).Return([]loadbalancer.Certificate{}, errors.New("test error")).Times(1)

		err := loadbalancerListenerCertificateList(service, &cobra.Command{}, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving certificates: test error", err.Error())
	})
}

func Test_loadbalancerListenerCertificateShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerListenerCertificateShowCmd(nil).Args(nil, []string{"123", "456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerCertificateShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing listener", err.Error())
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerCertificateShowCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing certificate", err.Error())
	})
}

func Test_loadbalancerListenerCertificateShow(t *testing.T) {
	t.Run("SingleListenerCertificate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetListenerCertificate(123, 456).Return(loadbalancer.Certificate{}, nil).Times(1)

		loadbalancerListenerCertificateShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("MultipleListenerCertificates", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetListenerCertificate(123, 456).Return(loadbalancer.Certificate{}, nil),
			service.EXPECT().GetListenerCertificate(123, 789).Return(loadbalancer.Certificate{}, nil),
		)

		loadbalancerListenerCertificateShow(service, &cobra.Command{}, []string{"123", "456", "789"})
	})

	t.Run("GetListenerID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		err := loadbalancerListenerCertificateShow(service, &cobra.Command{}, []string{"abc", "456"})

		assert.Equal(t, "Invalid listener ID", err.Error())
	})

	t.Run("GetCertificateID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid certificate ID [abc]\n", func() {
			loadbalancerListenerCertificateShow(service, &cobra.Command{}, []string{"123", "abc"})
		})
	})

	t.Run("GetListenerCertificateError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetListenerCertificate(123, 456).Return(loadbalancer.Certificate{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving certificate [456]: test error\n", func() {
			loadbalancerListenerCertificateShow(service, &cobra.Command{}, []string{"123", "456"})
		})
	})
}

func Test_loadbalancerListenerCertificateCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerListenerCertificateCreateCmd(nil, nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerCertificateCreateCmd(nil, nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing listener", err.Error())
	})
}

func Test_loadbalancerListenerCertificateCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerCertificateCreateCmd(nil, nil)
		cmd.ParseFlags([]string{"--name=test"})

		req := loadbalancer.CreateCertificateRequest{
			Name: "test",
		}

		gomock.InOrder(
			service.EXPECT().CreateListenerCertificate(123, req).Return(456, nil),
			service.EXPECT().GetListenerCertificate(123, 456).Return(loadbalancer.Certificate{}, nil),
		)

		loadbalancerListenerCertificateCreate(service, cmd, nil, []string{"123"})
	})

	t.Run("CreateListenerError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerCertificateCreateCmd(nil, nil)
		cmd.ParseFlags([]string{"--name=test"})

		service.EXPECT().CreateListenerCertificate(123, gomock.Any()).Return(0, errors.New("test error"))

		err := loadbalancerListenerCertificateCreate(service, cmd, nil, []string{"123"})

		assert.Equal(t, "Error creating certificate: test error", err.Error())
	})

	t.Run("GetListenerError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerCertificateCreateCmd(nil, nil)
		cmd.ParseFlags([]string{"--name=test"})

		gomock.InOrder(
			service.EXPECT().CreateListenerCertificate(123, gomock.Any()).Return(456, nil),
			service.EXPECT().GetListenerCertificate(123, 456).Return(loadbalancer.Certificate{}, errors.New("test error")),
		)

		err := loadbalancerListenerCertificateCreate(service, cmd, nil, []string{"123"})

		assert.Equal(t, "Error retrieving new certificate: test error", err.Error())
	})
}

func Test_loadbalancerListenerCertificateUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerListenerCertificateUpdateCmd(nil, nil).Args(nil, []string{"123", "456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerCertificateUpdateCmd(nil, nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing listener", err.Error())
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerCertificateUpdateCmd(nil, nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing certificate", err.Error())
	})
}

func Test_loadbalancerListenerCertificateUpdate(t *testing.T) {
	t.Run("SingleCertificate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		cmd := loadbalancerListenerCertificateUpdateCmd(nil, nil)
		cmd.ParseFlags([]string{"--name=test"})

		req := loadbalancer.PatchCertificateRequest{
			Name: "test",
		}

		gomock.InOrder(
			service.EXPECT().PatchListenerCertificate(123, 456, req).Return(nil),
			service.EXPECT().GetListenerCertificate(123, 456).Return(loadbalancer.Certificate{}, nil),
		)

		loadbalancerListenerCertificateUpdate(service, cmd, nil, []string{"123", "456"})
	})

	t.Run("PatchListenerError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().PatchListenerCertificate(123, 456, gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating certificate [456]: test error\n", func() {
			loadbalancerListenerCertificateUpdate(service, &cobra.Command{}, nil, []string{"123", "456"})
		})
	})

	t.Run("GetListenerError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchListenerCertificate(123, 456, gomock.Any()).Return(nil),
			service.EXPECT().GetListenerCertificate(123, 456).Return(loadbalancer.Certificate{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated certificate [456]: test error\n", func() {
			loadbalancerListenerCertificateUpdate(service, &cobra.Command{}, nil, []string{"123", "456"})
		})
	})
}

func Test_loadbalancerListenerCertificateDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerListenerCertificateDeleteCmd(nil).Args(nil, []string{"123", "456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerCertificateUpdateCmd(nil, nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing listener", err.Error())
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerCertificateUpdateCmd(nil, nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing certificate", err.Error())
	})
}

func Test_loadbalancerListenerCertificateDelete(t *testing.T) {
	t.Run("SingleCertificate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteListenerCertificate(123, 456).Return(nil).Times(1)

		loadbalancerListenerCertificateDelete(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("MultipleCertificates", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteListenerCertificate(1, 12).Return(nil),
			service.EXPECT().DeleteListenerCertificate(1, 123).Return(nil),
		)

		loadbalancerListenerCertificateDelete(service, &cobra.Command{}, []string{"1", "12", "123"})
	})

	t.Run("DeleteListenerError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteListenerCertificate(123, 456).Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing certificate [456]: test error\n", func() {
			loadbalancerListenerCertificateDelete(service, &cobra.Command{}, []string{"123", "456"})
		})
	})
}
