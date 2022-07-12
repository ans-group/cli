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

func Test_ecloudAffinityRuleList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetAffinityRules(gomock.Any()).Return([]ecloud.AffinityRule{}, nil).Times(1)

		ecloudAffinityRuleList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudAffinityRuleList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetAffinityRulesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetAffinityRules(gomock.Any()).Return([]ecloud.AffinityRule{}, errors.New("test error")).Times(1)

		err := ecloudAffinityRuleList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving affinity rules: test error", err.Error())
	})
}

func Test_ecloudAffinityRuleShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudAffinityRuleShowCmd(nil).Args(nil, []string{"ar-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudAffinityRuleShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing affinity rule", err.Error())
	})
}

func Test_ecloudAffinityRuleShow(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetAffinityRule("ar-abcdef12").Return(ecloud.AffinityRule{}, nil).Times(1)

		ecloudAffinityRuleShow(service, &cobra.Command{}, []string{"ar-abcdef12"})
	})

	t.Run("MultipleAffinityRules", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetAffinityRule("ar-abcdef12").Return(ecloud.AffinityRule{}, nil),
			service.EXPECT().GetAffinityRule("ar-abcdef23").Return(ecloud.AffinityRule{}, nil),
		)

		ecloudAffinityRuleShow(service, &cobra.Command{}, []string{"ar-abcdef12", "ar-abcdef23"})
	})

	t.Run("GetAffinityRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetAffinityRule("ar-abcdef12").Return(ecloud.AffinityRule{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving affinity rule [ar-abcdef12]: test error\n", func() {
			ecloudAffinityRuleShow(service, &cobra.Command{}, []string{"ar-abcdef12"})
		})
	})
}

func Test_ecloudAffinityRuleCreate(t *testing.T) {
	t.Run("CreateWithRequiredArgsNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudAffinityRuleCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--vpc=vpc-abcdef12", "--availability-zone=az-abcdef12", "--type=anti-affinity"})

		req := ecloud.CreateAffinityRuleRequest{
			Name: "testrule",
			VPCID: "vpc-abcdef12",
			AvailabilityZoneID: "az-abcdef12",
			Type: "anti-affinity",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "ar-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateAffinityRule(req).Return(resp, nil),
			service.EXPECT().GetAffinityRule("ar-abcdef12").Return(ecloud.AffinityRule{}, nil),
		)

		ecloudAffinityRuleCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudAffinityRuleCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--vpc=vpc-abcdef12", "--availability-zone=az-abcdef12", "--type=anti-affinity", "--wait"})

		req := ecloud.CreateAffinityRuleRequest{
			Name: "testrule",
			VPCID: "vpc-abcdef12",
			AvailabilityZoneID: "az-abcdef12",
			Type: "anti-affinity",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "ar-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateAffinityRule(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetAffinityRule("ar-abcdef12").Return(ecloud.AffinityRule{}, nil),
		)

		ecloudAffinityRuleCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudAffinityRuleCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--vpc=vpc-abcdef12", "--availability-zone=az-abcdef12", "--type=anti-affinity", "--wait"})

		req := ecloud.CreateAffinityRuleRequest{
			Name: "testrule",
			VPCID: "vpc-abcdef12",
			AvailabilityZoneID: "az-abcdef12",
			Type: "anti-affinity",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "ar-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateAffinityRule(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudAffinityRuleCreate(service, cmd, []string{})
		assert.Equal(t, "Error waiting for affinity rule task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateAffinityRuleError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudAffinityRuleCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--vpc=vpc-abcdef12", "--availability-zone=az-abcdef12", "--type=anti-affinity"})

		service.EXPECT().CreateAffinityRule(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudAffinityRuleCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating affinity rule: test error", err.Error())
	})

	t.Run("GetAffinityRuleError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudAffinityRuleCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--vpc=vpc-abcdef12", "--availability-zone=az-abcdef12", "--type=anti-affinity"})

		gomock.InOrder(
			service.EXPECT().CreateAffinityRule(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "ar-abcdef12"}, nil),
			service.EXPECT().GetAffinityRule("ar-abcdef12").Return(ecloud.AffinityRule{}, errors.New("test error")),
		)

		err := ecloudAffinityRuleCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new affinity rule: test error", err.Error())
	})
}

func Test_ecloudAffinityRuleUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudAffinityRuleUpdateCmd(nil).Args(nil, []string{"ar-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudAffinityRuleUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing affinity rule", err.Error())
	})
}

func Test_ecloudAffinityRuleUpdate(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudAffinityRuleUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule"})

		req := ecloud.PatchAffinityRuleRequest{
			Name: "testrule",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "ar-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchAffinityRule("ar-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetAffinityRule("ar-abcdef12").Return(ecloud.AffinityRule{}, nil),
		)

		ecloudAffinityRuleUpdate(service, cmd, []string{"ar-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudAffinityRuleUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--wait"})

		req := ecloud.PatchAffinityRuleRequest{
			Name: "testrule",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "ar-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchAffinityRule("ar-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetAffinityRule("ar-abcdef12").Return(ecloud.AffinityRule{}, nil),
		)

		ecloudAffinityRuleUpdate(service, cmd, []string{"ar-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudAffinityRuleUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--wait"})

		req := ecloud.PatchAffinityRuleRequest{
			Name: "testrule",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "ar-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchAffinityRule("ar-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for affinity rule [ar-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudAffinityRuleUpdate(service, cmd, []string{"ar-abcdef12"})
		})
	})

	t.Run("PatchAffinityRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchAffinityRule("ar-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating affinity rule [ar-abcdef12]: test error\n", func() {
			ecloudAffinityRuleUpdate(service, &cobra.Command{}, []string{"ar-abcdef12"})
		})
	})

	t.Run("GetAffinityRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "ar-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchAffinityRule("ar-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetAffinityRule("ar-abcdef12").Return(ecloud.AffinityRule{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated affinity rule [ar-abcdef12]: test error\n", func() {
			ecloudAffinityRuleUpdate(service, &cobra.Command{}, []string{"ar-abcdef12"})
		})
	})
}

func Test_ecloudAffinityRuleDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudAffinityRuleDeleteCmd(nil).Args(nil, []string{"ar-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudAffinityRuleDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing affinity rule", err.Error())
	})
}

func Test_ecloudAffinityRuleDelete(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteAffinityRule("ar-abcdef12").Return("task-abcdef12", nil)

		ecloudAffinityRuleDelete(service, &cobra.Command{}, []string{"ar-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudAffinityRuleDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteAffinityRule("ar-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudAffinityRuleDelete(service, cmd, []string{"ar-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudAffinityRuleDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteAffinityRule("ar-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for affinity rule [ar-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudAffinityRuleDelete(service, cmd, []string{"ar-abcdef12"})
		})
	})

	t.Run("DeleteAffinityRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteAffinityRule("ar-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing affinity rule [ar-abcdef12]: test error\n", func() {
			ecloudAffinityRuleDelete(service, &cobra.Command{}, []string{"ar-abcdef12"})
		})
	})
}
