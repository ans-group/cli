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

func Test_ecloudTemplateDiskListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudTemplateDiskListCmd().Args(nil, []string{"test template 1"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudTemplateDiskListCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_ecloudTemplateDiskList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTemplate("test template 1").Return(ecloud.Template{}, nil).Times(1)

		ecloudTemplateDiskList(service, &cobra.Command{}, []string{"test template 1"})
	})

	t.Run("GetTemplateError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTemplate("test template 1").Return(ecloud.Template{}, errors.New("test error 1")).Times(1)

		output := test.CatchStdErr(t, func() {
			ecloudTemplateDiskList(service, &cobra.Command{}, []string{"test template 1"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving template [test template 1]: test error 1\n", output)
	})
}
