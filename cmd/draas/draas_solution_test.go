package draas

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/draas"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_draasSolutionList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutions(gomock.Any()).Return([]draas.Solution{}, nil).Times(1)

		draasSolutionList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := draasSolutionList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSolutionsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutions(gomock.Any()).Return([]draas.Solution{}, errors.New("test error")).Times(1)

		err := draasSolutionList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving solutions: test error", err.Error())
	})
}

func Test_draasSolutionShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := draasSolutionShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := draasSolutionShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing solution", err.Error())
	})
}

func Test_draasSolutionShow(t *testing.T) {
	t.Run("SingleSolution", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolution("00000000-0000-0000-0000-000000000000").Return(draas.Solution{}, nil).Times(1)

		draasSolutionShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("GetSolutionError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolution("00000000-0000-0000-0000-000000000000").Return(draas.Solution{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving solution [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			draasSolutionShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_draasSolutionUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := draasSolutionUpdateCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := draasSolutionUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing solution", err.Error())
	})
}

func Test_draasSolutionUpdate(t *testing.T) {
	t.Run("DefaultUpdate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)
		cmd := draasSolutionUpdateCmd(nil)
		cmd.Flags().Set("name", "testname")

		req := draas.PatchSolutionRequest{
			Name: "testname",
		}

		gomock.InOrder(
			service.EXPECT().PatchSolution("00000000-0000-0000-0000-000000000000", req).Return(nil),
			service.EXPECT().GetSolution("00000000-0000-0000-0000-000000000000").Return(draas.Solution{}, nil),
		)

		draasSolutionUpdate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("PatchSolutionError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)
		cmd := draasSolutionUpdateCmd(nil)

		service.EXPECT().PatchSolution("00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error")).Times(1)

		err := draasSolutionUpdate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Equal(t, "error updating solution: test error", err.Error())
	})

	t.Run("GetSolutionError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)
		cmd := draasSolutionUpdateCmd(nil)

		gomock.InOrder(
			service.EXPECT().PatchSolution("00000000-0000-0000-0000-000000000000", gomock.Any()).Return(nil),
			service.EXPECT().GetSolution("00000000-0000-0000-0000-000000000000").Return(draas.Solution{}, errors.New("test error")),
		)

		err := draasSolutionUpdate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Equal(t, "error retrieving updated solution: test error", err.Error())
	})
}
