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

func Test_ecloudRouterCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudRouterCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrouter"})

		req := ecloud.CreateRouterRequest{
			Name: "testrouter",
		}

		gomock.InOrder(
			service.EXPECT().CreateRouter(req).Return("rtr-abcdef12", nil),
			service.EXPECT().GetRouter("rtr-abcdef12").Return(ecloud.Router{}, nil),
		)

		ecloudRouterCreate(service, cmd, []string{})
	})

	t.Run("CreateRouterError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudRouterCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrouter"})

		service.EXPECT().CreateRouter(gomock.Any()).Return("", errors.New("test error")).Times(1)

		err := ecloudRouterCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating router: test error", err.Error())
	})

	t.Run("GetRouterError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudRouterCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrouter"})

		gomock.InOrder(
			service.EXPECT().CreateRouter(gomock.Any()).Return("rtr-abcdef12", nil),
			service.EXPECT().GetRouter("rtr-abcdef12").Return(ecloud.Router{}, errors.New("test error")),
		)

		err := ecloudRouterCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new router: test error", err.Error())
	})
}

func Test_ecloudRouterUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudRouterUpdateCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudRouterUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing router", err.Error())
	})
}

func Test_ecloudRouterUpdate(t *testing.T) {
	t.Run("SingleRouter", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudRouterCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrouter"})

		req := ecloud.PatchRouterRequest{
			Name: "testrouter",
		}

		gomock.InOrder(
			service.EXPECT().PatchRouter("rtr-abcdef12", req).Return(nil),
			service.EXPECT().GetRouter("rtr-abcdef12").Return(ecloud.Router{}, nil),
		)

		ecloudRouterUpdate(service, cmd, []string{"rtr-abcdef12"})
	})

	t.Run("MultipleRouters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchRouter("rtr-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetRouter("rtr-abcdef12").Return(ecloud.Router{}, nil),
			service.EXPECT().PatchRouter("rtr-12abcdef", gomock.Any()).Return(nil),
			service.EXPECT().GetRouter("rtr-12abcdef").Return(ecloud.Router{}, nil),
		)

		ecloudRouterUpdate(service, &cobra.Command{}, []string{"rtr-abcdef12", "rtr-12abcdef"})
	})

	t.Run("PatchRouterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchRouter("rtr-abcdef12", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating router [rtr-abcdef12]: test error\n", func() {
			ecloudRouterUpdate(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})

	t.Run("GetRouterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchRouter("rtr-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetRouter("rtr-abcdef12").Return(ecloud.Router{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated router [rtr-abcdef12]: test error\n", func() {
			ecloudRouterUpdate(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})
}

func Test_ecloudRouterDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudRouterDeleteCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudRouterDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing router", err.Error())
	})
}

func Test_ecloudRouterDelete(t *testing.T) {
	t.Run("SingleRouter", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteRouter("rtr-abcdef12").Return(nil).Times(1)

		ecloudRouterDelete(service, &cobra.Command{}, []string{"rtr-abcdef12"})
	})

	t.Run("MultipleRouters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteRouter("rtr-abcdef12").Return(nil),
			service.EXPECT().DeleteRouter("rtr-12abcdef").Return(nil),
		)

		ecloudRouterDelete(service, &cobra.Command{}, []string{"rtr-abcdef12", "rtr-12abcdef"})
	})

	t.Run("DeleteRouterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteRouter("rtr-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing router [rtr-abcdef12]: test error\n", func() {
			ecloudRouterDelete(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})
}

func Test_ecloudRouterDeployDefaultFirewallPoliciesCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudRouterDeployDefaultFirewallPoliciesCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudRouterDeployDefaultFirewallPoliciesCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing router", err.Error())
	})
}

func Test_ecloudRouterDeployDefaultFirewallPolicies(t *testing.T) {
	t.Run("SingleRouter", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeployRouterDefaultFirewallPolicies("rtr-abcdef12").Return(nil)

		ecloudRouterDeployDefaultFirewallPolicies(service, &cobra.Command{}, []string{"rtr-abcdef12"})
	})

	t.Run("MultipleRouters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeployRouterDefaultFirewallPolicies("rtr-abcdef12").Return(nil),
			service.EXPECT().DeployRouterDefaultFirewallPolicies("rtr-abcdef23").Return(nil),
		)

		ecloudRouterDeployDefaultFirewallPolicies(service, &cobra.Command{}, []string{"rtr-abcdef12", "rtr-abcdef23"})
	})

	t.Run("DeployRouterDefaultFirewallPoliciesError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeployRouterDefaultFirewallPolicies("rtr-abcdef12").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error deploying default firewall policies for router [rtr-abcdef12]: test error\n", func() {
			ecloudRouterDeployDefaultFirewallPolicies(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})
}
