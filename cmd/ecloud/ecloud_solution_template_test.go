package ecloud

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudSolutionTemplateListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionTemplateListCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudSolutionTemplateListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})
}

func Test_ecloudSolutionTemplateList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutionTemplates(123, gomock.Any()).Return([]ecloud.Template{}, nil).Times(1)

		ecloudSolutionTemplateList(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("InvalidSolutionID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudSolutionTemplateList(service, &cobra.Command{}, []string{"abc"})

		assert.Equal(t, "Invalid solution ID [abc]", err.Error())
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudSolutionTemplateList(service, cmd, []string{"123"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetTemplatesError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutionTemplates(123, gomock.Any()).Return([]ecloud.Template{}, errors.New("test error 1")).Times(1)

		err := ecloudSolutionTemplateList(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error retrieving solution templates: test error 1", err.Error())
	})
}

func Test_ecloudSolutionTemplateShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionTemplateShowCmd(nil).Args(nil, []string{"123", "testtemplate1"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		err := ecloudSolutionTemplateShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})

	t.Run("MissingTemplate_Error", func(t *testing.T) {
		err := ecloudSolutionTemplateShowCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_ecloudSolutionTemplateShow(t *testing.T) {
	t.Run("RetrieveSingle", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutionTemplate(123, "testtemplate1").Return(ecloud.Template{}, nil).Times(1)

		ecloudSolutionTemplateShow(service, &cobra.Command{}, []string{"123", "testtemplate1"})
	})

	t.Run("RetrieveMultiple", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetSolutionTemplate(123, "testtemplate1").Return(ecloud.Template{}, nil),
			service.EXPECT().GetSolutionTemplate(123, "testtemplate2").Return(ecloud.Template{}, nil),
		)

		ecloudSolutionTemplateShow(service, &cobra.Command{}, []string{"123", "testtemplate1", "testtemplate2"})
	})

	t.Run("InvalidSolutionID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudSolutionTemplateShow(service, &cobra.Command{}, []string{"abc", "testtemplate1"})

		assert.Equal(t, "Invalid solution ID [abc]", err.Error())
	})

	t.Run("GetSolutionTemplateError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutionTemplate(123, "testtemplate1").Return(ecloud.Template{}, errors.New("test error 1")).Times(1)

		test_output.AssertErrorOutput(t, "Error retrieving solution template [testtemplate1]: test error 1\n", func() {
			ecloudSolutionTemplateShow(service, &cobra.Command{}, []string{"123", "testtemplate1"})
		})
	})
}

func Test_ecloudSolutionTemplateUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionTemplateUpdateCmd(nil).Args(nil, []string{"123", "testname1"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		err := ecloudSolutionTemplateUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})

	t.Run("MissingTemplate_Error", func(t *testing.T) {
		err := ecloudSolutionTemplateUpdateCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_ecloudSolutionTemplateUpdate(t *testing.T) {
	t.Run("UpdateSingle", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudSolutionTemplateUpdateCmd(nil)
		cmd.Flags().Set("name", "newname")

		expectedPatch := ecloud.RenameTemplateRequest{
			Destination: "newname",
		}

		gomock.InOrder(
			service.EXPECT().RenameSolutionTemplate(123, "testname1", gomock.Eq(expectedPatch)).Return(nil),
			service.EXPECT().GetSolutionTemplate(123, "newname").Return(ecloud.Template{}, nil),
			service.EXPECT().GetSolutionTemplate(123, "newname").Return(ecloud.Template{}, nil),
		)

		ecloudSolutionTemplateUpdate(service, cmd, []string{"123", "testname1"})
	})

	t.Run("InvalidSolutionID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudSolutionTemplateUpdate(service, &cobra.Command{}, []string{"abc"})

		assert.Equal(t, "Invalid solution ID [abc]", err.Error())
	})

	t.Run("RenameSolutionTemplateError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudSolutionTemplateUpdateCmd(nil)
		cmd.Flags().Set("name", "newname")

		gomock.InOrder(
			service.EXPECT().RenameSolutionTemplate(123, "testname1", gomock.Any()).Return(errors.New("test error 1")),
		)

		err := ecloudSolutionTemplateUpdate(service, cmd, []string{"123", "testname1"})

		assert.Equal(t, "Error updating solution template: test error 1", err.Error())
	})

	t.Run("WaitForCommandError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudSolutionTemplateUpdateCmd(nil)
		cmd.Flags().Set("name", "newname")

		gomock.InOrder(
			service.EXPECT().RenameSolutionTemplate(123, "testname1", gomock.Any()).Return(nil),
			service.EXPECT().GetSolutionTemplate(123, "newname").Return(ecloud.Template{}, errors.New("test error 1")),
		)

		err := ecloudSolutionTemplateUpdate(service, cmd, []string{"123", "testname1"})

		assert.Equal(t, "Error waiting for solution template update: Error waiting for command: Failed to retrieve solution template [newname]: test error 1", err.Error())
	})

	t.Run("GetSolutionTemplateError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetSolutionTemplate(123, "testname1").Return(ecloud.Template{}, errors.New("test error 1")),
		)

		err := ecloudSolutionTemplateUpdate(service, &cobra.Command{}, []string{"123", "testname1"})

		assert.Equal(t, "Error retrieving updated solution template: test error 1", err.Error())
	})
}

func Test_ecloudSolutionTemplateDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionTemplateDeleteCmd(nil).Args(nil, []string{"123", "testname1"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		err := ecloudSolutionTemplateDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})

	t.Run("MissingTemplate_Error", func(t *testing.T) {
		err := ecloudSolutionTemplateDeleteCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_ecloudSolutionTemplateDelete(t *testing.T) {
	t.Run("DeleteSingle", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteSolutionTemplate(123, "testname1").Return(nil).Times(1)

		ecloudSolutionTemplateDelete(service, &cobra.Command{}, []string{"123", "testname1"})
	})

	t.Run("DeleteMultiple", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteSolutionTemplate(123, "testname1").Return(nil),
			service.EXPECT().DeleteSolutionTemplate(123, "testname2").Return(nil),
		)

		ecloudSolutionTemplateDelete(service, &cobra.Command{}, []string{"123", "testname1", "testname2"})
	})

	t.Run("WithWait", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudSolutionTemplateDeleteCmd(nil)
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().DeleteSolutionTemplate(123, "testname1").Return(nil),
			service.EXPECT().GetSolutionTemplate(123, "testname1").Return(ecloud.Template{}, &ecloud.TemplateNotFoundError{}),
		)

		ecloudSolutionTemplateDelete(service, cmd, []string{"123", "testname1"})
	})

	t.Run("InvalidSolutionID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudSolutionTemplateDelete(service, &cobra.Command{}, []string{"abc", "testname1"})

		assert.Equal(t, "Invalid solution ID [abc]", err.Error())
	})

	t.Run("DeleteSolutionTemplateError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteSolutionTemplate(123, "testname1").Return(errors.New("test error 1")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing solution template [testname1]: test error 1\n", func() {
			ecloudSolutionTemplateDelete(service, &cobra.Command{}, []string{"123", "testname1"})
		})
	})

	t.Run("WaitForCommandError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudSolutionTemplateDeleteCmd(nil)
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().DeleteSolutionTemplate(123, "testname1").Return(nil),
			service.EXPECT().GetSolutionTemplate(123, "testname1").Return(ecloud.Template{}, errors.New("test error 1")),
		)

		test_output.AssertErrorOutput(t, "Error removing solution template [testname1]: Error waiting for command: Failed to retrieve solution template [testname1]: test error 1\n", func() {
			ecloudSolutionTemplateDelete(service, cmd, []string{"123", "testname1"})
		})
	})
}
