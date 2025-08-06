package pss

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_pssRequestFeedbackShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssRequestFeedbackShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := pssRequestFeedbackShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing request", err.Error())
	})
}

func Test_pssRequestFeedbackShow(t *testing.T) {
	t.Run("SingleRequest", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetRequestFeedback(123).Return(pss.Feedback{}, nil).Times(1)

		pssRequestFeedbackShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleRequests", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetRequestFeedback(123).Return(pss.Feedback{}, nil),
			service.EXPECT().GetRequestFeedback(456).Return(pss.Feedback{}, nil),
		)

		pssRequestFeedbackShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetRequestID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid request ID [abc]\n", func() {
			pssRequestFeedbackShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetRequestError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetRequestFeedback(123).Return(pss.Feedback{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving feedback for request [123]: test error\n", func() {
			pssRequestFeedbackShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_pssRequestFeedbackCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssRequestFeedbackCreateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := pssRequestFeedbackCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing request", err.Error())
	})
}

func Test_pssRequestFeedbackCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestFeedbackCreateCmd(nil)
		cmd.ParseFlags([]string{"--contact=1"})

		expectedRequest := pss.CreateFeedbackRequest{
			ContactID: 1,
		}

		gomock.InOrder(
			service.EXPECT().CreateRequestFeedback(123, gomock.Eq(expectedRequest)).Return(123, nil),
			service.EXPECT().GetRequestFeedback(123).Return(pss.Feedback{}, nil),
		)

		pssRequestFeedbackCreate(service, cmd, []string{"123"})
	})

	t.Run("InvalidRequestID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestFeedbackCreateCmd(nil)

		err := pssRequestFeedbackCreate(service, cmd, []string{"invalid"})
		assert.Contains(t, err.Error(), "invalid request ID [invalid]")
	})

	t.Run("CreateRequestFeedbackError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestFeedbackCreateCmd(nil)

		service.EXPECT().CreateRequestFeedback(123, gomock.Any()).Return(0, errors.New("test error")).Times(1)

		err := pssRequestFeedbackCreate(service, cmd, []string{"123"})
		assert.Equal(t, "error creating feedback for request: test error", err.Error())
	})

	t.Run("GetRequestError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestFeedbackCreateCmd(nil)

		gomock.InOrder(
			service.EXPECT().CreateRequestFeedback(123, gomock.Any()).Return(123, nil),
			service.EXPECT().GetRequestFeedback(123).Return(pss.Feedback{}, errors.New("test error")),
		)

		err := pssRequestFeedbackCreate(service, cmd, []string{"123"})
		assert.Equal(t, "error retrieving new feedback for request: test error", err.Error())
	})
}
