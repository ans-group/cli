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

func Test_pssChangeList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetChangeCases(gomock.Any()).Return([]pss.ChangeCase{}, nil).Times(1)

		pssChangeList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := pssChangeList(service, cmd, []string{})
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetChangesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetChangeCases(gomock.Any()).Return([]pss.ChangeCase{}, errors.New("test error")).Times(1)

		err := pssChangeList(service, &cobra.Command{}, []string{})
		assert.Equal(t, "test error", err.Error())
	})
}

// func Test_pssChangeShowCmd_Args(t *testing.T) {
// 	t.Run("ValidArgs_NoError", func(t *testing.T) {
// 		err := pssChangeShowCmd(nil).Args(nil, []string{"123"})

// 		assert.Nil(t, err)
// 	})

// 	t.Run("InvalidArgs_Error", func(t *testing.T) {
// 		err := pssChangeShowCmd(nil).Args(nil, []string{})

// 		assert.NotNil(t, err)
// 		assert.Equal(t, "Missing change", err.Error())
// 	})
// }

// func Test_pssChangeShow(t *testing.T) {
// 	t.Run("SingleChange", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)

// 		service.EXPECT().GetChange(123).Return(pss.Change{}, nil).Times(1)

// 		pssChangeShow(service, &cobra.Command{}, []string{"123"})
// 	})

// 	t.Run("MultipleChanges", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)

// 		gomock.InOrder(
// 			service.EXPECT().GetChange(123).Return(pss.Change{}, nil),
// 			service.EXPECT().GetChange(456).Return(pss.Change{}, nil),
// 		)

// 		pssChangeShow(service, &cobra.Command{}, []string{"123", "456"})
// 	})

// 	t.Run("GetChangeID_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)

// 		test_output.AssertErrorOutput(t, "Invalid change ID [abc]\n", func() {
// 			pssChangeShow(service, &cobra.Command{}, []string{"abc"})
// 		})
// 	})

// 	t.Run("GetChangeError_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)

// 		service.EXPECT().GetChange(123).Return(pss.Change{}, errors.New("test error"))

// 		test_output.AssertErrorOutput(t, "Error retrieving change [123]: test error\n", func() {
// 			pssChangeShow(service, &cobra.Command{}, []string{"123"})
// 		})
// 	})
// }

// func Test_pssChangeCreate(t *testing.T) {
// 	t.Run("DefaultCreate", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssChangeCreateCmd(nil)
// 		cmd.Flags().Set("subject", "test subject")
// 		cmd.Flags().Set("product-id", "456")
// 		cmd.Flags().Set("product-name", "testname")
// 		cmd.Flags().Set("product-type", "testtype")

// 		gomock.InOrder(
// 			service.EXPECT().CreateChange(gomock.Any()).Do(func(req pss.CreateChangeChange) {
// 				assert.Equal(t, "test subject", req.Subject)
// 				assert.Equal(t, 456, req.Product.ID)
// 				assert.Equal(t, "testname", req.Product.Name)
// 				assert.Equal(t, "testtype", req.Product.Type)
// 			}).Return(123, nil),
// 			service.EXPECT().GetChange(123).Return(pss.Change{}, nil),
// 		)

// 		pssChangeCreate(service, cmd, []string{})
// 	})

// 	t.Run("InvalidPriority_ReturnsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssChangeCreateCmd(nil)
// 		cmd.Flags().Set("priority", "invalid")

// 		err := pssChangeCreate(service, cmd, []string{})
// 		assert.Contains(t, err.Error(), "Invalid pss.ChangePriority")
// 	})

// 	t.Run("CreateChangeError_ReturnsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssChangeCreateCmd(nil)

// 		service.EXPECT().CreateChange(gomock.Any()).Return(0, errors.New("test error")).Times(1)

// 		err := pssChangeCreate(service, cmd, []string{})
// 		assert.Equal(t, "Error creating change: test error", err.Error())
// 	})

// 	t.Run("GetChangeError_ReturnsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssChangeCreateCmd(nil)

// 		gomock.InOrder(
// 			service.EXPECT().CreateChange(gomock.Any()).Return(123, nil),
// 			service.EXPECT().GetChange(123).Return(pss.Change{}, errors.New("test error")),
// 		)

// 		err := pssChangeCreate(service, cmd, []string{})
// 		assert.Equal(t, "Error retrieving new change: test error", err.Error())
// 	})
// }

// func Test_pssChangeUpdateCmd_Args(t *testing.T) {
// 	t.Run("ValidArgs_NoError", func(t *testing.T) {
// 		err := pssChangeUpdateCmd(nil).Args(nil, []string{"123"})

// 		assert.Nil(t, err)
// 	})

// 	t.Run("InvalidArgs_Error", func(t *testing.T) {
// 		err := pssChangeUpdateCmd(nil).Args(nil, []string{})

// 		assert.NotNil(t, err)
// 		assert.Equal(t, "Missing change", err.Error())
// 	})
// }

// func Test_pssChangeUpdate(t *testing.T) {
// 	t.Run("DefaultUpdate", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssChangeUpdateCmd(nil)
// 		cmd.Flags().Set("secure", "true")
// 		cmd.Flags().Set("read", "true")
// 		cmd.Flags().Set("change-sms", "true")
// 		cmd.Flags().Set("archived", "true")
// 		cmd.Flags().Set("priority", "High")

// 		gomock.InOrder(
// 			service.EXPECT().PatchChange(123, gomock.Any()).Do(func(changeID int, req pss.PatchChangeChange) {
// 				assert.Equal(t, 123, changeID)
// 				assert.Equal(t, true, *req.Secure)
// 				assert.Equal(t, true, *req.Read)
// 				assert.Equal(t, true, *req.ChangeSMS)
// 				assert.Equal(t, true, *req.Archived)
// 				assert.Equal(t, pss.ChangePriorityHigh, req.Priority)
// 			}).Return(nil),
// 			service.EXPECT().GetChange(123).Return(pss.Change{}, nil),
// 		)

// 		pssChangeUpdate(service, cmd, []string{"123"})
// 	})

// 	t.Run("InvalidPriority_ReturnsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssChangeUpdateCmd(nil)
// 		cmd.Flags().Set("priority", "invalid")

// 		err := pssChangeUpdate(service, cmd, []string{"123"})
// 		assert.Contains(t, err.Error(), "Invalid pss.ChangePriority")
// 	})

// 	t.Run("InvalidChangeID_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssChangeUpdateCmd(nil)

// 		test_output.AssertErrorOutput(t, "Invalid change ID [abc]\n", func() {
// 			pssChangeUpdate(service, cmd, []string{"abc"})
// 		})
// 	})

// 	t.Run("PatchChangeError_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssChangeUpdateCmd(nil)

// 		service.EXPECT().PatchChange(123, gomock.Any()).Return(errors.New("test error")).Times(1)

// 		test_output.AssertErrorOutput(t, "Error updating change [123]: test error\n", func() {
// 			pssChangeUpdate(service, cmd, []string{"123"})
// 		})
// 	})

// 	t.Run("GetChangeError_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssChangeUpdateCmd(nil)

// 		gomock.InOrder(
// 			service.EXPECT().PatchChange(123, gomock.Any()).Return(nil),
// 			service.EXPECT().GetChange(123).Return(pss.Change{}, errors.New("test error")),
// 		)

// 		test_output.AssertErrorOutput(t, "Error retrieving updated change [123]: test error\n", func() {
// 			pssChangeUpdate(service, cmd, []string{"123"})
// 		})
// 	})
// }

// func Test_pssChangeCloseCmd_Args(t *testing.T) {
// 	t.Run("ValidArgs_NoError", func(t *testing.T) {
// 		err := pssChangeCloseCmd(nil).Args(nil, []string{"123"})

// 		assert.Nil(t, err)
// 	})

// 	t.Run("InvalidArgs_Error", func(t *testing.T) {
// 		err := pssChangeCloseCmd(nil).Args(nil, []string{})

// 		assert.NotNil(t, err)
// 		assert.Equal(t, "Missing change", err.Error())
// 	})
// }

// func Test_pssChangeClose(t *testing.T) {
// 	t.Run("DefaultClose", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)

// 		gomock.InOrder(
// 			service.EXPECT().PatchChange(123, gomock.Any()).Return(nil),
// 			service.EXPECT().GetChange(123).Return(pss.Change{}, nil),
// 		)

// 		pssChangeClose(service, pssChangeCloseCmd(nil), []string{"123"})
// 	})

// 	t.Run("PatchChangeError_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssChangeCloseCmd(nil)

// 		service.EXPECT().PatchChange(123, gomock.Any()).Return(errors.New("test error")).Times(1)

// 		test_output.AssertErrorOutput(t, "Error closing change [123]: test error\n", func() {
// 			pssChangeClose(service, cmd, []string{"123"})
// 		})
// 	})

// 	t.Run("GetChangeError_OutputsError", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()

// 		service := mocks.NewMockPSSService(mockCtrl)
// 		cmd := pssChangeCloseCmd(nil)

// 		gomock.InOrder(
// 			service.EXPECT().PatchChange(123, gomock.Any()).Return(nil),
// 			service.EXPECT().GetChange(123).Return(pss.Change{}, errors.New("test error")),
// 		)

// 		test_output.AssertErrorOutput(t, "Error retrieving updated change [123]: test error\n", func() {
// 			pssChangeClose(service, cmd, []string{"123"})
// 		})
// 	})
// }
