package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudSolutionTemplateListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionTemplateListCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudSolutionTemplateListCmd().Args(nil, []string{})

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

	t.Run("InvalidSolutionID_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		output := test.CatchStdErr(t, func() {
			ecloudSolutionTemplateList(service, &cobra.Command{}, []string{"abc"})
		})

		assert.Equal(t, "Invalid solution ID [abc]\n", output)
		assert.Equal(t, 1, code)
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
			ecloudSolutionTemplateList(service, &cobra.Command{}, []string{"123"})
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

		service.EXPECT().GetSolutionTemplates(123, gomock.Any()).Return([]ecloud.Template{}, errors.New("test error 1")).Times(1)

		output := test.CatchStdErr(t, func() {
			ecloudSolutionTemplateList(service, &cobra.Command{}, []string{"123"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving solution templates: test error 1\n", output)
	})
}

func Test_ecloudSolutionTemplateShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionTemplateShowCmd().Args(nil, []string{"123", "testtemplate1"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		err := ecloudSolutionTemplateShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})

	t.Run("MissingTemplate_Error", func(t *testing.T) {
		err := ecloudSolutionTemplateShowCmd().Args(nil, []string{"123"})

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

	t.Run("InvalidSolutionID_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		output := test.CatchStdErr(t, func() {
			ecloudSolutionTemplateShow(service, &cobra.Command{}, []string{"abc", "testtemplate1"})
		})

		assert.Equal(t, "Invalid solution ID [abc]\n", output)
		assert.Equal(t, 1, code)
	})

	t.Run("GetSolutionTemplateError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutionTemplate(123, "testtemplate1").Return(ecloud.Template{}, errors.New("test error 1")).Times(1)

		output := test.CatchStdErr(t, func() {
			ecloudSolutionTemplateShow(service, &cobra.Command{}, []string{"123", "testtemplate1"})
		})

		assert.Equal(t, "Error retrieving solution template [testtemplate1]: test error 1\n", output)
	})
}
