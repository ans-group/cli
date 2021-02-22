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

func Test_ecloudFirewallPolicyList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallPolicies(gomock.Any()).Return([]ecloud.FirewallPolicy{}, nil).Times(1)

		ecloudFirewallPolicyList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudFirewallPolicyList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetFirewallPoliciesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallPolicies(gomock.Any()).Return([]ecloud.FirewallPolicy{}, errors.New("test error")).Times(1)

		err := ecloudFirewallPolicyList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving firewall policies: test error", err.Error())
	})
}

func Test_ecloudFirewallPolicyShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFirewallPolicyShowCmd(nil).Args(nil, []string{"fwp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFirewallPolicyShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing firewall policy", err.Error())
	})
}

func Test_ecloudFirewallPolicyShow(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, nil).Times(1)

		ecloudFirewallPolicyShow(service, &cobra.Command{}, []string{"fwp-abcdef12"})
	})

	t.Run("MultipleFirewallPolicies", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, nil),
			service.EXPECT().GetFirewallPolicy("fwp-abcdef23").Return(ecloud.FirewallPolicy{}, nil),
		)

		ecloudFirewallPolicyShow(service, &cobra.Command{}, []string{"fwp-abcdef12", "fwp-abcdef23"})
	})

	t.Run("GetFirewallPolicyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving firewall policy [fwp-abcdef12]: test error\n", func() {
			ecloudFirewallPolicyShow(service, &cobra.Command{}, []string{"fwp-abcdef12"})
		})
	})
}

func Test_ecloudFirewallPolicyCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFirewallPolicyCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		req := ecloud.CreateFirewallPolicyRequest{
			Name: "testpolicy",
		}

		gomock.InOrder(
			service.EXPECT().CreateFirewallPolicy(req).Return("fwp-abcdef12", nil),
			service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, nil),
		)

		ecloudFirewallPolicyCreate(service, cmd, []string{})
	})

	t.Run("CreateFirewallPolicyError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFirewallPolicyCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		service.EXPECT().CreateFirewallPolicy(gomock.Any()).Return("", errors.New("test error")).Times(1)

		err := ecloudFirewallPolicyCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating firewall policy: test error", err.Error())
	})

	t.Run("GetFirewallPolicyError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFirewallPolicyCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		gomock.InOrder(
			service.EXPECT().CreateFirewallPolicy(gomock.Any()).Return("fwp-abcdef12", nil),
			service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, errors.New("test error")),
		)

		err := ecloudFirewallPolicyCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new firewall policy: test error", err.Error())
	})
}

func Test_ecloudFirewallPolicyUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFirewallPolicyUpdateCmd(nil).Args(nil, []string{"fwp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFirewallPolicyUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing firewall policy", err.Error())
	})
}

func Test_ecloudFirewallPolicyUpdate(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudFirewallPolicyUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		req := ecloud.PatchFirewallPolicyRequest{
			Name: "testpolicy",
		}

		gomock.InOrder(
			service.EXPECT().PatchFirewallPolicy("fwp-abcdef12", req).Return(nil),
			service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, nil),
		)

		ecloudFirewallPolicyUpdate(service, cmd, []string{"fwp-abcdef12"})
	})

	t.Run("MultipleFirewallPolicies", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchFirewallPolicy("fwp-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, nil),
			service.EXPECT().PatchFirewallPolicy("fwp-12abcdef", gomock.Any()).Return(nil),
			service.EXPECT().GetFirewallPolicy("fwp-12abcdef").Return(ecloud.FirewallPolicy{}, nil),
		)

		ecloudFirewallPolicyUpdate(service, &cobra.Command{}, []string{"fwp-abcdef12", "fwp-12abcdef"})
	})

	t.Run("PatchFirewallPolicyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchFirewallPolicy("fwp-abcdef12", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating firewall policy [fwp-abcdef12]: test error\n", func() {
			ecloudFirewallPolicyUpdate(service, &cobra.Command{}, []string{"fwp-abcdef12"})
		})
	})

	t.Run("GetFirewallPolicyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchFirewallPolicy("fwp-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated firewall policy [fwp-abcdef12]: test error\n", func() {
			ecloudFirewallPolicyUpdate(service, &cobra.Command{}, []string{"fwp-abcdef12"})
		})
	})
}

func Test_ecloudFirewallPolicyDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFirewallPolicyDeleteCmd(nil).Args(nil, []string{"fwp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFirewallPolicyDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing firewall policy", err.Error())
	})
}

func Test_ecloudFirewallPolicyDelete(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteFirewallPolicy("fwp-abcdef12").Return(nil).Times(1)

		ecloudFirewallPolicyDelete(service, &cobra.Command{}, []string{"fwp-abcdef12"})
	})

	t.Run("MultipleFirewallPolicies", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteFirewallPolicy("fwp-abcdef12").Return(nil),
			service.EXPECT().DeleteFirewallPolicy("fwp-12abcdef").Return(nil),
		)

		ecloudFirewallPolicyDelete(service, &cobra.Command{}, []string{"fwp-abcdef12", "fwp-12abcdef"})
	})

	t.Run("DeleteFirewallPolicyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteFirewallPolicy("fwp-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing firewall policy [fwp-abcdef12]: test error\n", func() {
			ecloudFirewallPolicyDelete(service, &cobra.Command{}, []string{"fwp-abcdef12"})
		})
	})
}
