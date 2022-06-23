package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudInstanceConsoleSessionCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudInstanceConsoleSessionCreateCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("MissingInstance_Error", func(t *testing.T) {
		err := ecloudInstanceConsoleSessionCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing instance", err.Error())
	})
}

func Test_ecloudInstanceConsoleSessionCreate(t *testing.T) {
	t.Run("CreateSuccess_NoError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().CreateInstanceConsoleSession("i-abcdef12").Return(ecloud.ConsoleSession{}, nil).Times(1)

		err := ecloudInstanceConsoleSessionCreate(service, &cobra.Command{}, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("GetInstanceConsoleSessionError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().CreateInstanceConsoleSession("i-abcdef12").Return(ecloud.ConsoleSession{}, errors.New("test error 1")).Times(1)

		test_output.AssertErrorOutput(t, "Error creating instance [i-abcdef12] console session: test error 1\n", func() {
			ecloudInstanceConsoleSessionCreate(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})
}
