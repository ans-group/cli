package pss

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_pssProblemList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetProblemCases(gomock.Any()).Return([]pss.ProblemCase{}, nil).Times(1)

		pssProblemList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := pssProblemList(service, cmd, []string{})
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetProblemsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetProblemCases(gomock.Any()).Return([]pss.ProblemCase{}, errors.New("test error")).Times(1)

		err := pssProblemList(service, &cobra.Command{}, []string{})
		assert.Equal(t, "test error", err.Error())
	})
}

func Test_pssProblemShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssProblemShowCmd(nil).Args(nil, []string{"PRB123456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := pssProblemShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing problem", err.Error())
	})
}

func Test_pssProblemShow(t *testing.T) {
	t.Run("SingleProblem", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetProblemCase("PRB123456").Return(pss.ProblemCase{}, nil)

		pssProblemShow(service, &cobra.Command{}, []string{"PRB123456"})
	})

	t.Run("MultipleProblems", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetProblemCase("PRB123").Return(pss.ProblemCase{}, nil),
			service.EXPECT().GetProblemCase("PRB456").Return(pss.ProblemCase{}, nil),
		)

		pssProblemShow(service, &cobra.Command{}, []string{"PRB123", "PRB456"})
	})

	t.Run("GetProblemError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetProblemCase("PRB123456").Return(pss.ProblemCase{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving problem [PRB123456]: test error\n", func() {
			pssProblemShow(service, &cobra.Command{}, []string{"PRB123456"})
		})
	})
}
