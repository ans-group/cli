package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func Test_loadtestTestJobCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadtestTestJobCreateCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadtestTestJobCreateCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing test", err.Error())
	})
}

func Test_loadtestTestJobCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)
		cmd := loadtestTestJobCreateCmd()
		cmd.Flags().Set("run-now", "true")

		expectedRequest := ltaas.CreateTestJobRequest{
			RunNow: true,
		}

		gomock.InOrder(
			service.EXPECT().CreateTestJob("00000000-0000-0000-0000-000000000000", expectedRequest).Return("00000000-0000-0000-0000-000000000001", nil),
			service.EXPECT().GetJob("00000000-0000-0000-0000-000000000001").Return(ltaas.Job{}, nil),
		)

		loadtestTestJobCreate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("CreateJobError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)
		cmd := loadtestTestJobCreateCmd()

		service.EXPECT().CreateTestJob("00000000-0000-0000-0000-000000000000", gomock.Any()).Return("", errors.New("test error")).Times(1)

		err := loadtestTestJobCreate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})
		assert.Equal(t, "Error creating job: test error", err.Error())
	})

	t.Run("GetJobError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)
		cmd := loadtestTestJobCreateCmd()

		gomock.InOrder(
			service.EXPECT().CreateTestJob("00000000-0000-0000-0000-000000000000", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetJob("00000000-0000-0000-0000-000000000000").Return(ltaas.Job{}, errors.New("test error")),
		)

		err := loadtestTestJobCreate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})
		assert.Equal(t, "Error retrieving new job [00000000-0000-0000-0000-000000000000]: test error", err.Error())
	})
}
