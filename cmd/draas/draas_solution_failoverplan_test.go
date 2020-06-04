package draas

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/draas"
)

func Test_draasSolutionFailoverPlanListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := draasSolutionFailoverPlanListCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		cmd := draasSolutionFailoverPlanListCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})
}

func Test_draasSolutionFailoverPlanList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionFailoverPlans("00000000-0000-0000-0000-000000000000", gomock.Any()).Return([]draas.FailoverPlan{}, nil).Times(1)

		draasSolutionFailoverPlanList(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := draasSolutionFailoverPlanList(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSolutionFailoverPlansError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionFailoverPlans("00000000-0000-0000-0000-000000000000", gomock.Any()).Return([]draas.FailoverPlan{}, errors.New("test error")).Times(1)

		err := draasSolutionFailoverPlanList(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Equal(t, "Error retrieving solution failover plans: test error", err.Error())
	})
}

func Test_draasSolutionFailoverPlanShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := draasSolutionFailoverPlanShowCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		cmd := draasSolutionFailoverPlanShowCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})

	t.Run("MissingFailoverPlan_Error", func(t *testing.T) {
		cmd := draasSolutionFailoverPlanShowCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing failover plan", err.Error())
	})
}

func Test_draasSolutionFailoverPlanShow(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionFailoverPlan("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001").Return(draas.FailoverPlan{}, nil)

		draasSolutionFailoverPlanShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("GetSolutionFailoverPlan_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionFailoverPlan("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001").Return(draas.FailoverPlan{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving solution failover plan [00000000-0000-0000-0000-000000000001]: test error\n", func() {
			draasSolutionFailoverPlanShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
		})
	})
}

func Test_draasSolutionFailoverPlanStartCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := draasSolutionFailoverPlanStartCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		cmd := draasSolutionFailoverPlanStartCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})

	t.Run("MissingFailoverPlan_Error", func(t *testing.T) {
		cmd := draasSolutionFailoverPlanStartCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing failover plan", err.Error())
	})
}

func Test_draasSolutionFailoverPlanStart(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		req := draas.StartFailoverPlanRequest{}

		service.EXPECT().StartSolutionFailoverPlan("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001", req).Return(nil)

		draasSolutionFailoverPlanStart(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("StartSolutionFailoverPlanError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().StartSolutionFailoverPlan("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error starting solution failover plan [00000000-0000-0000-0000-000000000001]: test error\n", func() {
			draasSolutionFailoverPlanStart(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
		})
	})
}

func Test_draasSolutionFailoverPlanStopCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := draasSolutionFailoverPlanStopCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		cmd := draasSolutionFailoverPlanStopCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})

	t.Run("MissingFailoverPlan_Error", func(t *testing.T) {
		cmd := draasSolutionFailoverPlanStopCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing failover plan", err.Error())
	})
}

func Test_draasSolutionFailoverPlanStop(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().StopSolutionFailoverPlan("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001").Return(nil)

		draasSolutionFailoverPlanStop(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("StopSolutionFailoverPlan_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().StopSolutionFailoverPlan("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error stopping solution failover plan [00000000-0000-0000-0000-000000000001]: test error\n", func() {
			draasSolutionFailoverPlanStop(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
		})
	})
}
