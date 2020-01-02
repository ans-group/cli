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

func Test_ecloudSolutionHostListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionHostListCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudSolutionHostListCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})
}

func Test_ecloudSolutionHostList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutionHosts(123, gomock.Any()).Return([]ecloud.Host{}, nil).Times(1)

		ecloudSolutionHostList(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("InvalidSolutionID_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		test_output.AssertFatalOutput(t, "Invalid solution ID [abc]\n", func() {
			ecloudSolutionHostList(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		test_output.AssertFatalOutput(t, "Missing value for filtering\n", func() {
			ecloudSolutionHostList(service, cmd, []string{"123"})
		})
	})

	t.Run("GetHostsError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutionHosts(123, gomock.Any()).Return([]ecloud.Host{}, errors.New("test error 1")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving solution hosts: test error 1\n", func() {
			ecloudSolutionHostList(service, &cobra.Command{}, []string{"123"})
		})
	})
}
