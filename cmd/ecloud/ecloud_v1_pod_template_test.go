package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/config"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudPodTemplateListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudPodTemplateListCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudPodTemplateListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing pod", err.Error())
	})
}

func Test_ecloudPodTemplateList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetPodTemplates(123, gomock.Any()).Return([]ecloud.Template{}, nil).Times(1)

		ecloudPodTemplateList(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("InvalidPodID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudPodTemplateList(service, &cobra.Command{}, []string{"abc"})

		assert.Equal(t, "Invalid pod ID [abc]", err.Error())
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudPodTemplateList(service, cmd, []string{"123"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetTemplatesError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetPodTemplates(123, gomock.Any()).Return([]ecloud.Template{}, errors.New("test error 1")).Times(1)

		err := ecloudPodTemplateList(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error retrieving pod templates: test error 1", err.Error())
	})
}

func Test_ecloudPodTemplateShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudPodTemplateShowCmd(nil).Args(nil, []string{"123", "testtemplate1"})

		assert.Nil(t, err)
	})

	t.Run("MissingPod_Error", func(t *testing.T) {
		err := ecloudPodTemplateShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing pod", err.Error())
	})

	t.Run("MissingTemplate_Error", func(t *testing.T) {
		err := ecloudPodTemplateShowCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_ecloudPodTemplateShow(t *testing.T) {
	t.Run("RetrieveSingle", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetPodTemplate(123, "testtemplate1").Return(ecloud.Template{}, nil).Times(1)

		ecloudPodTemplateShow(service, &cobra.Command{}, []string{"123", "testtemplate1"})
	})

	t.Run("RetrieveMultiple", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetPodTemplate(123, "testtemplate1").Return(ecloud.Template{}, nil),
			service.EXPECT().GetPodTemplate(123, "testtemplate2").Return(ecloud.Template{}, nil),
		)

		ecloudPodTemplateShow(service, &cobra.Command{}, []string{"123", "testtemplate1", "testtemplate2"})
	})

	t.Run("InvalidPodID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudPodTemplateShow(service, &cobra.Command{}, []string{"abc", "testtemplate1"})

		assert.Equal(t, "Invalid pod ID [abc]", err.Error())
	})

	t.Run("GetPodTemplateError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetPodTemplate(123, "testtemplate1").Return(ecloud.Template{}, errors.New("test error 1")).Times(1)

		test_output.AssertErrorOutput(t, "Error retrieving pod template [testtemplate1]: test error 1\n", func() {
			ecloudPodTemplateShow(service, &cobra.Command{}, []string{"123", "testtemplate1"})
		})
	})
}

func Test_ecloudPodTemplateUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudPodTemplateUpdateCmd(nil).Args(nil, []string{"123", "testname1"})

		assert.Nil(t, err)
	})

	t.Run("MissingPod_Error", func(t *testing.T) {
		err := ecloudPodTemplateUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing pod", err.Error())
	})

	t.Run("MissingTemplate_Error", func(t *testing.T) {
		err := ecloudPodTemplateUpdateCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_ecloudPodTemplateUpdate(t *testing.T) {
	t.Run("UpdateSingle", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		config.Set("test", "command_wait_timeout_seconds", 1200)
		config.Set("test", "command_wait_sleep_seconds", 1)
		config.SwitchCurrentContext("test")
		defer config.Reset()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudPodTemplateUpdateCmd(nil)
		cmd.Flags().Set("name", "newname")

		expectedPatch := ecloud.RenameTemplateRequest{
			Destination: "newname",
		}

		gomock.InOrder(
			service.EXPECT().RenamePodTemplate(123, "testname1", gomock.Eq(expectedPatch)).Return(nil),
			service.EXPECT().GetPodTemplate(123, "newname").Return(ecloud.Template{}, nil),
			service.EXPECT().GetPodTemplate(123, "newname").Return(ecloud.Template{}, nil),
		)

		ecloudPodTemplateUpdate(service, cmd, []string{"123", "testname1"})
	})

	t.Run("InvalidPodID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudPodTemplateUpdate(service, &cobra.Command{}, []string{"abc"})

		assert.Equal(t, "Invalid pod ID [abc]", err.Error())
	})

	t.Run("RenamePodTemplateError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudPodTemplateUpdateCmd(nil)
		cmd.Flags().Set("name", "newname")

		gomock.InOrder(
			service.EXPECT().RenamePodTemplate(123, "testname1", gomock.Any()).Return(errors.New("test error 1")),
		)

		err := ecloudPodTemplateUpdate(service, cmd, []string{"123", "testname1"})

		assert.Equal(t, "Error updating pod template: test error 1", err.Error())
	})

	t.Run("WaitForCommandError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		config.Set("test", "command_wait_timeout_seconds", 1200)
		config.Set("test", "command_wait_sleep_seconds", 1)
		config.SwitchCurrentContext("test")
		defer config.Reset()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudPodTemplateUpdateCmd(nil)
		cmd.Flags().Set("name", "newname")

		gomock.InOrder(
			service.EXPECT().RenamePodTemplate(123, "testname1", gomock.Any()).Return(nil),
			service.EXPECT().GetPodTemplate(123, "newname").Return(ecloud.Template{}, errors.New("test error 1")),
		)

		err := ecloudPodTemplateUpdate(service, cmd, []string{"123", "testname1"})

		assert.Equal(t, "Error waiting for pod template update: Error waiting for command: Failed to retrieve pod template [newname]: test error 1", err.Error())
	})

	t.Run("GetPodTemplateError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetPodTemplate(123, "testname1").Return(ecloud.Template{}, errors.New("test error 1")),
		)

		err := ecloudPodTemplateUpdate(service, &cobra.Command{}, []string{"123", "testname1"})

		assert.Equal(t, "Error retrieving updated pod template: test error 1", err.Error())
	})
}

func Test_ecloudPodTemplateDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudPodTemplateDeleteCmd(nil).Args(nil, []string{"123", "testname1"})

		assert.Nil(t, err)
	})

	t.Run("MissingPod_Error", func(t *testing.T) {
		err := ecloudPodTemplateDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing pod", err.Error())
	})

	t.Run("MissingTemplate_Error", func(t *testing.T) {
		err := ecloudPodTemplateDeleteCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_ecloudPodTemplateDelete(t *testing.T) {
	t.Run("DeleteSingle", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeletePodTemplate(123, "testname1").Return(nil).Times(1)

		ecloudPodTemplateDelete(service, &cobra.Command{}, []string{"123", "testname1"})
	})

	t.Run("DeleteMultiple", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeletePodTemplate(123, "testname1").Return(nil),
			service.EXPECT().DeletePodTemplate(123, "testname2").Return(nil),
		)

		ecloudPodTemplateDelete(service, &cobra.Command{}, []string{"123", "testname1", "testname2"})
	})

	t.Run("WithWait", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		config.Set("test", "command_wait_timeout_seconds", 1200)
		config.Set("test", "command_wait_sleep_seconds", 1)
		config.SwitchCurrentContext("test")
		defer config.Reset()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudPodTemplateDeleteCmd(nil)
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().DeletePodTemplate(123, "testname1").Return(nil),
			service.EXPECT().GetPodTemplate(123, "testname1").Return(ecloud.Template{}, &ecloud.TemplateNotFoundError{}),
		)

		ecloudPodTemplateDelete(service, cmd, []string{"123", "testname1"})
	})

	t.Run("InvalidPodID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudPodTemplateDelete(service, &cobra.Command{}, []string{"abc", "testname1"})

		assert.Equal(t, "Invalid pod ID [abc]", err.Error())
	})

	t.Run("DeletePodTemplateError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeletePodTemplate(123, "testname1").Return(errors.New("test error 1")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing pod template [testname1]: test error 1\n", func() {
			ecloudPodTemplateDelete(service, &cobra.Command{}, []string{"123", "testname1"})
		})
	})

	t.Run("WaitForCommandError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		config.Set("test", "command_wait_timeout_seconds", 1200)
		config.Set("test", "command_wait_sleep_seconds", 1)
		config.SwitchCurrentContext("test")
		defer config.Reset()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudPodTemplateDeleteCmd(nil)
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().DeletePodTemplate(123, "testname1").Return(nil),
			service.EXPECT().GetPodTemplate(123, "testname1").Return(ecloud.Template{}, errors.New("test error 1")),
		)

		test_output.AssertErrorOutput(t, "Error removing pod template [testname1]: Error waiting for command: Failed to retrieve pod template [testname1]: test error 1\n", func() {
			ecloudPodTemplateDelete(service, cmd, []string{"123", "testname1"})
		})
	})
}
