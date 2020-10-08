package loadbalancer

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func Test_loadbalancerConfigurationList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetConfigurations(gomock.Any()).Return([]loadbalancer.Configuration{}, nil).Times(1)

		loadbalancerConfigurationList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadbalancerConfigurationList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetConfigurationsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetConfigurations(gomock.Any()).Return([]loadbalancer.Configuration{}, errors.New("test error")).Times(1)

		err := loadbalancerConfigurationList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving configurations: test error", err.Error())
	})
}

func Test_loadbalancerConfigurationShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerConfigurationShowCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerConfigurationShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing configuration", err.Error())
	})
}

func Test_loadbalancerConfigurationShow(t *testing.T) {
	t.Run("SingleConfiguration", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetConfiguration("rtr-abcdef12").Return(loadbalancer.Configuration{}, nil).Times(1)

		loadbalancerConfigurationShow(service, &cobra.Command{}, []string{"rtr-abcdef12"})
	})

	t.Run("MultipleConfigurations", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetConfiguration("rtr-abcdef12").Return(loadbalancer.Configuration{}, nil),
			service.EXPECT().GetConfiguration("rtr-abcdef23").Return(loadbalancer.Configuration{}, nil),
		)

		loadbalancerConfigurationShow(service, &cobra.Command{}, []string{"rtr-abcdef12", "rtr-abcdef23"})
	})

	t.Run("GetConfigurationError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetConfiguration("rtr-abcdef12").Return(loadbalancer.Configuration{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving configuration [rtr-abcdef12]: test error\n", func() {
			loadbalancerConfigurationShow(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})
}

func Test_loadbalancerConfigurationCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerConfigurationCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testconfiguration", "--vpc=vpc-abcdef12", "--az=az-abcdef12"})

		req := loadbalancer.CreateConfigurationRequest{
			Name: ptr.String("testconfiguration"),
		}

		gomock.InOrder(
			service.EXPECT().CreateConfiguration(req).Return("rtr-abcdef12", nil),
			service.EXPECT().GetConfiguration("rtr-abcdef12").Return(loadbalancer.Configuration{}, nil),
		)

		loadbalancerConfigurationCreate(service, cmd, []string{})
	})

	t.Run("CreateConfigurationError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerConfigurationCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testconfiguration", "--configuration=rtr-abcdef12"})

		service.EXPECT().CreateConfiguration(gomock.Any()).Return("", errors.New("test error")).Times(1)

		err := loadbalancerConfigurationCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating configuration: test error", err.Error())
	})

	t.Run("GetConfigurationError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerConfigurationCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testconfiguration", "--configuration=rtr-abcdef12"})

		gomock.InOrder(
			service.EXPECT().CreateConfiguration(gomock.Any()).Return("rtr-abcdef12", nil),
			service.EXPECT().GetConfiguration("rtr-abcdef12").Return(loadbalancer.Configuration{}, errors.New("test error")),
		)

		err := loadbalancerConfigurationCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new configuration: test error", err.Error())
	})
}

func Test_loadbalancerConfigurationUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerConfigurationUpdateCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerConfigurationUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing configuration", err.Error())
	})
}

func Test_loadbalancerConfigurationUpdate(t *testing.T) {
	t.Run("SingleConfiguration", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		cmd := loadbalancerConfigurationCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testconfiguration"})

		req := loadbalancer.PatchConfigurationRequest{
			Name: ptr.String("testconfiguration"),
		}

		gomock.InOrder(
			service.EXPECT().PatchConfiguration("rtr-abcdef12", req).Return(nil),
			service.EXPECT().GetConfiguration("rtr-abcdef12").Return(loadbalancer.Configuration{}, nil),
		)

		loadbalancerConfigurationUpdate(service, cmd, []string{"rtr-abcdef12"})
	})

	t.Run("MultipleConfigurations", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchConfiguration("rtr-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetConfiguration("rtr-abcdef12").Return(loadbalancer.Configuration{}, nil),
			service.EXPECT().PatchConfiguration("rtr-12abcdef", gomock.Any()).Return(nil),
			service.EXPECT().GetConfiguration("rtr-12abcdef").Return(loadbalancer.Configuration{}, nil),
		)

		loadbalancerConfigurationUpdate(service, &cobra.Command{}, []string{"rtr-abcdef12", "rtr-12abcdef"})
	})

	t.Run("PatchConfigurationError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().PatchConfiguration("rtr-abcdef12", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating configuration [rtr-abcdef12]: test error\n", func() {
			loadbalancerConfigurationUpdate(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})

	t.Run("GetConfigurationError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchConfiguration("rtr-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetConfiguration("rtr-abcdef12").Return(loadbalancer.Configuration{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated configuration [rtr-abcdef12]: test error\n", func() {
			loadbalancerConfigurationUpdate(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})
}

func Test_loadbalancerConfigurationDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerConfigurationDeleteCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerConfigurationDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing configuration", err.Error())
	})
}

func Test_loadbalancerConfigurationDelete(t *testing.T) {
	t.Run("SingleConfiguration", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteConfiguration("rtr-abcdef12").Return(nil).Times(1)

		loadbalancerConfigurationDelete(service, &cobra.Command{}, []string{"rtr-abcdef12"})
	})

	t.Run("MultipleConfigurations", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteConfiguration("rtr-abcdef12").Return(nil),
			service.EXPECT().DeleteConfiguration("rtr-12abcdef").Return(nil),
		)

		loadbalancerConfigurationDelete(service, &cobra.Command{}, []string{"rtr-abcdef12", "rtr-12abcdef"})
	})

	t.Run("DeleteConfigurationError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteConfiguration("rtr-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing configuration [rtr-abcdef12]: test error\n", func() {
			loadbalancerConfigurationDelete(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})
}
