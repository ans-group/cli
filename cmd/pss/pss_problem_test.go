package pss

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
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

// func Test_pssProblemShowCmd_Args(t *testing.T) {
// 	t.Run("ValidArgs_NoError", func(t *testing.T) {
// 		err := pssProblemShowCmd(nil).Args(nil, []string{"123"})

// 		assert.Nil(t, err)
// 	})

// 	t.Run("InvalidArgs_Error", func(t *testing.T) {
// 		err := pssProblemShowCmd(nil).Args(nil, []string{})

// 		assert.NotNil(t, err)
// 		assert.Equal(t, "Missing problem", err.Error())
// 	})
// }

// func Test_pssProblemShow(t *testing.T) {
// 	t.Run("SingleProblem", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)

// 		service.EXPECT().GetProblem(123).Return(pss.Problem{}, nil).Times(1)

// 		pssProblemShow(service, &cobra.Command{}, []string{"123"})
// 	})

// 	t.Run("MultipleProblems", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)

// 		gomock.InOrder(
// 			service.EXPECT().GetProblem(123).Return(pss.Problem{}, nil),
// 			service.EXPECT().GetProblem(456).Return(pss.Problem{}, nil),
// 		)

// 		pssProblemShow(service, &cobra.Command{}, []string{"123", "456"})
// 	})

// 	t.Run("GetProblemID_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)

// 		test_output.AssertErrorOutput(t, "Invalid problem ID [abc]\n", func() {
// 			pssProblemShow(service, &cobra.Command{}, []string{"abc"})
// 		})
// 	})

// 	t.Run("GetProblemError_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)

// 		service.EXPECT().GetProblem(123).Return(pss.Problem{}, errors.New("test error"))

// 		test_output.AssertErrorOutput(t, "Error retrieving problem [123]: test error\n", func() {
// 			pssProblemShow(service, &cobra.Command{}, []string{"123"})
// 		})
// 	})
// }

// func Test_pssProblemCreate(t *testing.T) {
// 	t.Run("DefaultCreate", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssProblemCreateCmd(nil)
// 		cmd.Flags().Set("subject", "test subject")
// 		cmd.Flags().Set("product-id", "456")
// 		cmd.Flags().Set("product-name", "testname")
// 		cmd.Flags().Set("product-type", "testtype")

// 		gomock.InOrder(
// 			service.EXPECT().CreateProblem(gomock.Any()).Do(func(req pss.CreateProblemProblem) {
// 				assert.Equal(t, "test subject", req.Subject)
// 				assert.Equal(t, 456, req.Product.ID)
// 				assert.Equal(t, "testname", req.Product.Name)
// 				assert.Equal(t, "testtype", req.Product.Type)
// 			}).Return(123, nil),
// 			service.EXPECT().GetProblem(123).Return(pss.Problem{}, nil),
// 		)

// 		pssProblemCreate(service, cmd, []string{})
// 	})

// 	t.Run("InvalidPriority_ReturnsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssProblemCreateCmd(nil)
// 		cmd.Flags().Set("priority", "invalid")

// 		err := pssProblemCreate(service, cmd, []string{})
// 		assert.Contains(t, err.Error(), "Invalid pss.ProblemPriority")
// 	})

// 	t.Run("CreateProblemError_ReturnsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssProblemCreateCmd(nil)

// 		service.EXPECT().CreateProblem(gomock.Any()).Return(0, errors.New("test error")).Times(1)

// 		err := pssProblemCreate(service, cmd, []string{})
// 		assert.Equal(t, "Error creating problem: test error", err.Error())
// 	})

// 	t.Run("GetProblemError_ReturnsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssProblemCreateCmd(nil)

// 		gomock.InOrder(
// 			service.EXPECT().CreateProblem(gomock.Any()).Return(123, nil),
// 			service.EXPECT().GetProblem(123).Return(pss.Problem{}, errors.New("test error")),
// 		)

// 		err := pssProblemCreate(service, cmd, []string{})
// 		assert.Equal(t, "Error retrieving new problem: test error", err.Error())
// 	})
// }

// func Test_pssProblemUpdateCmd_Args(t *testing.T) {
// 	t.Run("ValidArgs_NoError", func(t *testing.T) {
// 		err := pssProblemUpdateCmd(nil).Args(nil, []string{"123"})

// 		assert.Nil(t, err)
// 	})

// 	t.Run("InvalidArgs_Error", func(t *testing.T) {
// 		err := pssProblemUpdateCmd(nil).Args(nil, []string{})

// 		assert.NotNil(t, err)
// 		assert.Equal(t, "Missing problem", err.Error())
// 	})
// }

// func Test_pssProblemUpdate(t *testing.T) {
// 	t.Run("DefaultUpdate", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssProblemUpdateCmd(nil)
// 		cmd.Flags().Set("secure", "true")
// 		cmd.Flags().Set("read", "true")
// 		cmd.Flags().Set("problem-sms", "true")
// 		cmd.Flags().Set("archived", "true")
// 		cmd.Flags().Set("priority", "High")

// 		gomock.InOrder(
// 			service.EXPECT().PatchProblem(123, gomock.Any()).Do(func(problemID int, req pss.PatchProblemProblem) {
// 				assert.Equal(t, 123, problemID)
// 				assert.Equal(t, true, *req.Secure)
// 				assert.Equal(t, true, *req.Read)
// 				assert.Equal(t, true, *req.ProblemSMS)
// 				assert.Equal(t, true, *req.Archived)
// 				assert.Equal(t, pss.ProblemPriorityHigh, req.Priority)
// 			}).Return(nil),
// 			service.EXPECT().GetProblem(123).Return(pss.Problem{}, nil),
// 		)

// 		pssProblemUpdate(service, cmd, []string{"123"})
// 	})

// 	t.Run("InvalidPriority_ReturnsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssProblemUpdateCmd(nil)
// 		cmd.Flags().Set("priority", "invalid")

// 		err := pssProblemUpdate(service, cmd, []string{"123"})
// 		assert.Contains(t, err.Error(), "Invalid pss.ProblemPriority")
// 	})

// 	t.Run("InvalidProblemID_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssProblemUpdateCmd(nil)

// 		test_output.AssertErrorOutput(t, "Invalid problem ID [abc]\n", func() {
// 			pssProblemUpdate(service, cmd, []string{"abc"})
// 		})
// 	})

// 	t.Run("PatchProblemError_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssProblemUpdateCmd(nil)

// 		service.EXPECT().PatchProblem(123, gomock.Any()).Return(errors.New("test error")).Times(1)

// 		test_output.AssertErrorOutput(t, "Error updating problem [123]: test error\n", func() {
// 			pssProblemUpdate(service, cmd, []string{"123"})
// 		})
// 	})

// 	t.Run("GetProblemError_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssProblemUpdateCmd(nil)

// 		gomock.InOrder(
// 			service.EXPECT().PatchProblem(123, gomock.Any()).Return(nil),
// 			service.EXPECT().GetProblem(123).Return(pss.Problem{}, errors.New("test error")),
// 		)

// 		test_output.AssertErrorOutput(t, "Error retrieving updated problem [123]: test error\n", func() {
// 			pssProblemUpdate(service, cmd, []string{"123"})
// 		})
// 	})
// }

// func Test_pssProblemCloseCmd_Args(t *testing.T) {
// 	t.Run("ValidArgs_NoError", func(t *testing.T) {
// 		err := pssProblemCloseCmd(nil).Args(nil, []string{"123"})

// 		assert.Nil(t, err)
// 	})

// 	t.Run("InvalidArgs_Error", func(t *testing.T) {
// 		err := pssProblemCloseCmd(nil).Args(nil, []string{})

// 		assert.NotNil(t, err)
// 		assert.Equal(t, "Missing problem", err.Error())
// 	})
// }

// func Test_pssProblemClose(t *testing.T) {
// 	t.Run("DefaultClose", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)

// 		gomock.InOrder(
// 			service.EXPECT().PatchProblem(123, gomock.Any()).Return(nil),
// 			service.EXPECT().GetProblem(123).Return(pss.Problem{}, nil),
// 		)

// 		pssProblemClose(service, pssProblemCloseCmd(nil), []string{"123"})
// 	})

// 	t.Run("PatchProblemError_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssProblemCloseCmd(nil)

// 		service.EXPECT().PatchProblem(123, gomock.Any()).Return(errors.New("test error")).Times(1)

// 		test_output.AssertErrorOutput(t, "Error closing problem [123]: test error\n", func() {
// 			pssProblemClose(service, cmd, []string{"123"})
// 		})
// 	})

// 	t.Run("GetProblemError_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssProblemCloseCmd(nil)

// 		gomock.InOrder(
// 			service.EXPECT().PatchProblem(123, gomock.Any()).Return(nil),
// 			service.EXPECT().GetProblem(123).Return(pss.Problem{}, errors.New("test error")),
// 		)

// 		test_output.AssertErrorOutput(t, "Error retrieving updated problem [123]: test error\n", func() {
// 			pssProblemClose(service, cmd, []string{"123"})
// 		})
// 	})
// }
