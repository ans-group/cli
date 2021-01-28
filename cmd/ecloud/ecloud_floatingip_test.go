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

func Test_ecloudFloatingIPCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFloatingIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testfip"})

		req := ecloud.CreateFloatingIPRequest{
			Name: "testfip",
		}

		gomock.InOrder(
			service.EXPECT().CreateFloatingIP(req).Return("fip-abcdef12", nil),
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil),
		)

		ecloudFloatingIPCreate(service, cmd, []string{})
	})

	t.Run("CreateFloatingIPError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFloatingIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testfip"})

		service.EXPECT().CreateFloatingIP(gomock.Any()).Return("", errors.New("test error")).Times(1)

		err := ecloudFloatingIPCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating floating IP: test error", err.Error())
	})

	t.Run("GetFloatingIPError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFloatingIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testfip"})

		gomock.InOrder(
			service.EXPECT().CreateFloatingIP(gomock.Any()).Return("fip-abcdef12", nil),
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, errors.New("test error")),
		)

		err := ecloudFloatingIPCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new floating IP: test error", err.Error())
	})
}

func Test_ecloudFloatingIPUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFloatingIPUpdateCmd(nil).Args(nil, []string{"fip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFloatingIPUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing floating IP", err.Error())
	})
}

func Test_ecloudFloatingIPUpdate(t *testing.T) {
	t.Run("SingleFloatingIP", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudFloatingIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testfip"})

		req := ecloud.PatchFloatingIPRequest{
			Name: "testfip",
		}

		gomock.InOrder(
			service.EXPECT().PatchFloatingIP("fip-abcdef12", req).Return(nil),
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil),
		)

		ecloudFloatingIPUpdate(service, cmd, []string{"fip-abcdef12"})
	})

	t.Run("MultipleFloatingIPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchFloatingIP("fip-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil),
			service.EXPECT().PatchFloatingIP("fip-12abcdef", gomock.Any()).Return(nil),
			service.EXPECT().GetFloatingIP("fip-12abcdef").Return(ecloud.FloatingIP{}, nil),
		)

		ecloudFloatingIPUpdate(service, &cobra.Command{}, []string{"fip-abcdef12", "fip-12abcdef"})
	})

	t.Run("PatchFloatingIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchFloatingIP("fip-abcdef12", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating floating IP [fip-abcdef12]: test error\n", func() {
			ecloudFloatingIPUpdate(service, &cobra.Command{}, []string{"fip-abcdef12"})
		})
	})

	t.Run("GetFloatingIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchFloatingIP("fip-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated floating IP [fip-abcdef12]: test error\n", func() {
			ecloudFloatingIPUpdate(service, &cobra.Command{}, []string{"fip-abcdef12"})
		})
	})
}

func Test_ecloudFloatingIPDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFloatingIPDeleteCmd(nil).Args(nil, []string{"fip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFloatingIPDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing floating IP", err.Error())
	})
}

func Test_ecloudFloatingIPDelete(t *testing.T) {
	t.Run("SingleFloatingIP", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteFloatingIP("fip-abcdef12").Return(nil).Times(1)

		ecloudFloatingIPDelete(service, &cobra.Command{}, []string{"fip-abcdef12"})
	})

	t.Run("MultipleFloatingIPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteFloatingIP("fip-abcdef12").Return(nil),
			service.EXPECT().DeleteFloatingIP("fip-12abcdef").Return(nil),
		)

		ecloudFloatingIPDelete(service, &cobra.Command{}, []string{"fip-abcdef12", "fip-12abcdef"})
	})

	t.Run("DeleteFloatingIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteFloatingIP("fip-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing floating IP [fip-abcdef12]: test error\n", func() {
			ecloudFloatingIPDelete(service, &cobra.Command{}, []string{"fip-abcdef12"})
		})
	})
}

func Test_ecloudFloatingIPAssignCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFloatingIPAssignCmd(nil).Args(nil, []string{"fip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFloatingIPAssignCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing floating IP", err.Error())
	})
}

func Test_ecloudFloatingIPAssign(t *testing.T) {
	t.Run("AssignFloatingIP_NoError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		req := ecloud.AssignFloatingIPRequest{
			ResourceID: "i-abcdef12",
		}

		cmd := ecloudFloatingIPAssignCmd(nil)
		cmd.ParseFlags([]string{"--resource=i-abcdef12"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().AssignFloatingIP("fip-abcdef12", gomock.Eq(req)).Return(nil)
		service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil)

		err := ecloudFloatingIPAssign(service, cmd, []string{"fip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("AssignFloatingIPError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().AssignFloatingIP("fip-abcdef12", gomock.Any()).Return(errors.New("test error"))

		err := ecloudFloatingIPAssign(service, &cobra.Command{}, []string{"fip-abcdef12"})

		assert.NotNil(t, err)
	})
}

func Test_ecloudFloatingIPUnassignCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFloatingIPUnassignCmd(nil).Args(nil, []string{"fip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFloatingIPUnassignCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing floating IP", err.Error())
	})
}

func Test_ecloudFloatingIPUnassign(t *testing.T) {
	t.Run("SingleFloatingIP", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().UnassignFloatingIP("fip-abcdef12").Return(nil).Times(1)

		ecloudFloatingIPUnassign(service, &cobra.Command{}, []string{"fip-abcdef12"})
	})

	t.Run("MultipleFloatingIPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().UnassignFloatingIP("fip-abcdef12").Return(nil),
			service.EXPECT().UnassignFloatingIP("fip-12abcdef").Return(nil),
		)

		ecloudFloatingIPUnassign(service, &cobra.Command{}, []string{"fip-abcdef12", "fip-12abcdef"})
	})

	t.Run("UnassignFloatingIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().UnassignFloatingIP("fip-abcdef12").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error unassigning floating IP [fip-abcdef12]: test error\n", func() {
			ecloudFloatingIPUnassign(service, &cobra.Command{}, []string{"fip-abcdef12"})
		})
	})
}
