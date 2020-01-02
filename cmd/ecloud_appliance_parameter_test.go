package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudApplianceParameterListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudApplianceParameterListCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudApplianceParameterListCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing appliance", err.Error())
	})
}

func Test_ecloudApplianceParameterList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetApplianceParameters("00000000-0000-0000-0000-000000000000", gomock.Any()).Return([]ecloud.ApplianceParameter{}, nil).Times(1)

		ecloudApplianceParameterList(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		test_output.AssertFatalOutput(t, "Missing value for filtering\n", func() {
			ecloudApplianceParameterList(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})

	t.Run("GetParametersError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetApplianceParameters("00000000-0000-0000-0000-000000000000", gomock.Any()).Return([]ecloud.ApplianceParameter{}, errors.New("test error 1")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving appliance parameters: test error 1\n", func() {
			ecloudApplianceParameterList(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}
