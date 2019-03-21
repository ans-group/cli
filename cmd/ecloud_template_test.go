package cmd

import (
	"errors"
	"testing"

	"github.com/ukfast/sdk-go/pkg/ptr"

	"github.com/ukfast/cli/internal/pkg/output"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudTemplateList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTemplates(gomock.Any()).Return([]ecloud.Template{}, nil).Times(1)

		ecloudTemplateList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		output := test.CatchStdErr(t, func() {
			ecloudTemplateList(service, &cobra.Command{}, []string{"test template 1"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Missing value for filtering\n", output)
	})

	t.Run("GetTemplatesError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})

		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTemplates(gomock.Any()).Return([]ecloud.Template{}, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			ecloudTemplateList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving templates: test error\n", output)
	})
}

func Test_ecloudTemplateShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudTemplateShowCmd().Args(nil, []string{"test template 1"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudTemplateShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_ecloudTemplateShow(t *testing.T) {
	t.Run("SingleTemplate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTemplate("test template 1").Return(ecloud.Template{}, nil).Times(1)

		ecloudTemplateShow(service, &cobra.Command{}, []string{"test template 1"})
	})

	t.Run("MultipleTemplates", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetTemplate("test template 1").Return(ecloud.Template{}, nil),
			service.EXPECT().GetTemplate("test template 2").Return(ecloud.Template{}, nil),
		)

		ecloudTemplateShow(service, &cobra.Command{}, []string{"test template 1", "test template 2"})
	})

	t.Run("GetTemplateError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTemplate("test template 1").Return(ecloud.Template{}, errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			ecloudTemplateShow(service, &cobra.Command{}, []string{"test template 1"})
		})

		assert.Equal(t, "Error retrieving template [test template 1]: test error\n", output)
	})
}

func Test_ecloudTemplateUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudTemplateUpdateCmd().Args(nil, []string{"test template 1"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudTemplateUpdateCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_ecloudTemplateUpdate(t *testing.T) {
	t.Run("Update", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudTemplateUpdateCmd()
		cmd.Flags().Set("name", "new template name")
		cmd.Flags().Set("solution", "12345")

		expectedRequest := ecloud.RenameTemplateRequest{
			NewTemplateName: "new template name",
			SolutionID:      ptr.Int(12345),
		}

		gomock.InOrder(
			service.EXPECT().RenameTemplate("test template 1", expectedRequest).Return(nil),
			service.EXPECT().GetTemplate("new template name").Return(ecloud.Template{}, nil),
		)

		ecloudTemplateUpdate(service, cmd, []string{"test template 1"})
	})

	t.Run("RenameTemplateError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().RenameTemplate("test template 1", gomock.Any()).Return(errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			ecloudTemplateUpdate(service, &cobra.Command{}, []string{"test template 1"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error updating template [test template 1]: test error\n", output)
	})

	t.Run("GetTemplateError_OutputsError", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudTemplateUpdateCmd()
		cmd.Flags().Set("name", "new template name")

		gomock.InOrder(
			service.EXPECT().RenameTemplate("test template 1", gomock.Any()).Return(nil),
			service.EXPECT().GetTemplate("new template name").Return(ecloud.Template{}, errors.New("test error")),
		)

		output := test.CatchStdErr(t, func() {
			ecloudTemplateUpdate(service, cmd, []string{"test template 1"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving updated template: test error\n", output)
	})
}

func Test_ecloudTemplateDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudTemplateDeleteCmd().Args(nil, []string{"test template 1"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudTemplateDeleteCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_ecloudTemplateDelete(t *testing.T) {
	t.Run("SingleTemplate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteTemplate("test template 1").Return(nil).Times(1)

		ecloudTemplateDelete(service, &cobra.Command{}, []string{"test template 1"})
	})

	t.Run("MultipleTemplates", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteTemplate("test template 1").Return(nil),
			service.EXPECT().DeleteTemplate("test template 2").Return(nil),
		)

		ecloudTemplateDelete(service, &cobra.Command{}, []string{"test template 1", "test template 2"})
	})

	t.Run("DeleteTemplateError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteTemplate("test template 1").Return(errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			ecloudTemplateDelete(service, &cobra.Command{}, []string{"test template 1"})
		})

		assert.Equal(t, "Error removing template [test template 1]: test error\n", output)
	})
}
