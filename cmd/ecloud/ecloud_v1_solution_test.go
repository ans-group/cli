package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudSolutionList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutions(gomock.Any()).Return([]ecloud.Solution{}, nil).Times(1)

		ecloudSolutionList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudSolutionList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSolutionsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolutions(gomock.Any()).Return([]ecloud.Solution{}, errors.New("test error")).Times(1)

		err := ecloudSolutionList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving solutions: test error", err.Error())
	})
}

func Test_ecloudSolutionShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudSolutionShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})
}

func Test_ecloudSolutionShow(t *testing.T) {
	t.Run("SingleSolution", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolution(123).Return(ecloud.Solution{}, nil).Times(1)

		ecloudSolutionShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleSolutions", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetSolution(123).Return(ecloud.Solution{}, nil),
			service.EXPECT().GetSolution(456).Return(ecloud.Solution{}, nil),
		)

		ecloudSolutionShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetSolutionID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid solution ID [abc]\n", func() {
			ecloudSolutionShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetSolutionError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetSolution(123).Return(ecloud.Solution{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving solution [123]: test error\n", func() {
			ecloudSolutionShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_ecloudSolutionUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudSolutionUpdateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudSolutionUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})
}

func Test_ecloudSolutionUpdate(t *testing.T) {
	t.Run("Update", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudSolutionUpdateCmd(nil)
		cmd.Flags().Set("name", "new solution name")

		gomock.InOrder(
			service.EXPECT().PatchSolution(123, gomock.Any()).Do(func(solutionID int, patch ecloud.PatchSolutionRequest) {
				if patch.Name == nil || *patch.Name != "new solution name" {
					t.Fatal("Unexpected solution name")
				}
			}).Return(123, nil),
			service.EXPECT().GetSolution(123).Return(ecloud.Solution{}, nil),
		)

		ecloudSolutionUpdate(service, cmd, []string{"123"})
	})

	t.Run("InvalidSolutionID_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudSolutionUpdate(service, &cobra.Command{}, []string{"abc"})

		assert.Equal(t, "Invalid solution ID [abc]", err.Error())
	})

	t.Run("PatchSolutionError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchSolution(123, gomock.Any()).Return(0, errors.New("test error")).Times(1)

		err := ecloudSolutionUpdate(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error updating solution: test error", err.Error())
	})

	t.Run("GetSolutionError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchSolution(123, gomock.Any()).Return(123, nil),
			service.EXPECT().GetSolution(123).Return(ecloud.Solution{}, errors.New("test error")),
		)

		err := ecloudSolutionUpdate(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error retrieving updated solution: test error", err.Error())
	})
}
