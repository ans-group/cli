package loadtest

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func Test_loadtestJobList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetJobs(gomock.Any()).Return([]ltaas.Job{}, nil).Times(1)

		loadtestJobList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadtestJobList(service, cmd, []string{})

		assert.NotNil(t, err)
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetJobsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetJobs(gomock.Any()).Return([]ltaas.Job{}, errors.New("test error")).Times(1)

		err := loadtestJobList(service, &cobra.Command{}, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving jobs: test error", err.Error())
	})
}

func Test_loadtestJobShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadtestJobShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadtestJobShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing job", err.Error())
	})
}

func Test_loadtestJobShow(t *testing.T) {
	t.Run("SingleJob", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetJob("00000000-0000-0000-0000-000000000000").Return(ltaas.Job{}, nil).Times(1)

		loadtestJobShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MultipleJobs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetJob("00000000-0000-0000-0000-000000000000").Return(ltaas.Job{}, nil),
			service.EXPECT().GetJob("00000000-0000-0000-0000-000000000001").Return(ltaas.Job{}, nil),
		)

		loadtestJobShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("GetJobError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetJob("00000000-0000-0000-0000-000000000000").Return(ltaas.Job{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving job [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			loadtestJobShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_loadtestJobCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)
		cmd := loadtestJobCreateCmd(nil)
		cmd.Flags().Set("test-id", "00000000-0000-0000-0000-000000000000")

		expectedRequest := ltaas.CreateJobRequest{
			TestID: "00000000-0000-0000-0000-000000000000",
		}

		gomock.InOrder(
			service.EXPECT().CreateJob(expectedRequest).Return("00000000-0000-0000-0000-000000000001", nil),
			service.EXPECT().GetJob("00000000-0000-0000-0000-000000000001").Return(ltaas.Job{}, nil),
		)

		loadtestJobCreate(service, cmd, []string{})
	})

	t.Run("CreateJobError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)
		cmd := loadtestJobCreateCmd(nil)

		service.EXPECT().CreateJob(gomock.Any()).Return("", errors.New("test error")).Times(1)

		err := loadtestJobCreate(service, cmd, []string{})
		assert.Equal(t, "Error creating job: test error", err.Error())
	})

	t.Run("GetJobError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)
		cmd := loadtestJobCreateCmd(nil)

		gomock.InOrder(
			service.EXPECT().CreateJob(gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetJob("00000000-0000-0000-0000-000000000000").Return(ltaas.Job{}, errors.New("test error")),
		)

		err := loadtestJobCreate(service, cmd, []string{})
		assert.Equal(t, "Error retrieving new job [00000000-0000-0000-0000-000000000000]: test error", err.Error())
	})
}

func Test_loadtestJobDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadtestJobDeleteCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadtestJobDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing job", err.Error())
	})
}

func Test_loadtestJobDelete(t *testing.T) {
	t.Run("SingleJob", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetJob("00000000-0000-0000-0000-000000000000").Return(ltaas.Job{}, nil).Times(1)

		loadtestJobDelete(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MultipleJobs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetJob("00000000-0000-0000-0000-000000000000").Return(ltaas.Job{}, nil),
			service.EXPECT().GetJob("00000000-0000-0000-0000-000000000001").Return(ltaas.Job{}, nil),
		)

		loadtestJobDelete(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("GetJobError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetJob("00000000-0000-0000-0000-000000000000").Return(ltaas.Job{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing job [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			loadtestJobDelete(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_loadtestJobStopCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadtestJobStopCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadtestJobStopCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing job", err.Error())
	})
}

func Test_loadtestJobStop(t *testing.T) {
	t.Run("SingleJob", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().StopJob("00000000-0000-0000-0000-000000000000").Return(nil).Times(1)

		loadtestJobStop(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MultipleJobs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().StopJob("00000000-0000-0000-0000-000000000000").Return(nil),
			service.EXPECT().StopJob("00000000-0000-0000-0000-000000000001").Return(nil),
		)

		loadtestJobStop(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("GetJobError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().StopJob("00000000-0000-0000-0000-000000000000").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error stopping job [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			loadtestJobStop(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}
