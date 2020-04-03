package loadtest

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func Test_loadtestJobSettingsShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadtestJobSettingsShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadtestJobSettingsShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing job", err.Error())
	})
}

func Test_loadtestJobSettingsShow(t *testing.T) {
	t.Run("SingleJob", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetJobSettings("00000000-0000-0000-0000-000000000000").Return(ltaas.JobSettings{}, nil).Times(1)

		loadtestJobSettingsShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MultipleJobs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetJobSettings("00000000-0000-0000-0000-000000000000").Return(ltaas.JobSettings{}, nil),
			service.EXPECT().GetJobSettings("00000000-0000-0000-0000-000000000001").Return(ltaas.JobSettings{}, nil),
		)

		loadtestJobSettingsShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("GetJobError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetJobSettings("00000000-0000-0000-0000-000000000000").Return(ltaas.JobSettings{}, errors.New("test error"))

		err := loadtestJobSettingsShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Equal(t, "Error retrieving job settings: test error", err.Error())
	})
}
