package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/storage"
)

func Test_storageSolutionList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		service.EXPECT().GetSolutions(gomock.Any()).Return([]storage.Solution{}, nil).Times(1)

		storageSolutionList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		test_output.AssertFatalOutput(t, "Missing value for filtering\n", func() {
			storageSolutionList(service, cmd, []string{})
		})
	})

	t.Run("GetSolutionsError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		service.EXPECT().GetSolutions(gomock.Any()).Return([]storage.Solution{}, errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving solutions: test error\n", func() {
			storageSolutionList(service, &cobra.Command{}, []string{})
		})
	})
}

func Test_storageSolutionShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := storageSolutionShowCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := storageSolutionShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})
}

func Test_storageSolutionShow(t *testing.T) {
	t.Run("SingleSolution", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		service.EXPECT().GetSolution(123).Return(storage.Solution{}, nil).Times(1)

		storageSolutionShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleSolutions", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetSolution(123).Return(storage.Solution{}, nil),
			service.EXPECT().GetSolution(456).Return(storage.Solution{}, nil),
		)

		storageSolutionShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetSolutionID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid solution ID [abc]\n", func() {
			storageSolutionShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetSolutionError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockStorageService(mockCtrl)

		service.EXPECT().GetSolution(123).Return(storage.Solution{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving solution [123]: test error\n", func() {
			storageSolutionShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
