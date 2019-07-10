package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/pss"
)

func Test_pssRequestReplyListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssRequestReplyListCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := pssRequestReplyListCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing request", err.Error())
	})
}

func Test_pssRequestReplyList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetRequestConversation(123, gomock.Any()).Return([]pss.Reply{}, nil).Times(1)

		pssRequestReplyList(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("GetRequestID_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		test_output.AssertFatalOutput(t, "Invalid request ID [abc]\n", func() {
			pssRequestReplyList(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		test_output.AssertFatalOutput(t, "Missing value for filtering\n", func() {
			pssRequestReplyList(service, &cobra.Command{}, []string{"123"})
		})
	})

	t.Run("GetRequestConversationError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetRequestConversation(123, gomock.Any()).Return([]pss.Reply{}, errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving request replies: test error\n", func() {
			pssRequestReplyList(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_pssRequestReplyCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssRequestReplyCreateCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := pssRequestReplyCreateCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing request", err.Error())
	})
}

func Test_pssRequestReplyCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestReplyCreateCmd()
		cmd.Flags().Set("description", "test description")
		cmd.Flags().Set("author", "456")

		expectedRequest := pss.CreateReplyRequest{
			Description: "test description",
			Author: pss.Author{
				ID: 456,
			},
		}

		gomock.InOrder(
			service.EXPECT().CreateRequestReply(123, gomock.Eq(expectedRequest)).Return("C456", nil),
			service.EXPECT().GetReply("C456").Return(pss.Reply{}, nil),
		)

		pssRequestReplyCreate(service, cmd, []string{"123"})
	})

	t.Run("InvalidRequestID_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestReplyCreateCmd()

		test_output.AssertFatalOutputFunc(t, func(stdErr string) {
			assert.Contains(t, stdErr, "Invalid request ID [invalid]")
		}, func() {
			pssRequestReplyCreate(service, cmd, []string{"invalid"})
		})
	})

	t.Run("CreateRequestReplyError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestReplyCreateCmd()

		service.EXPECT().CreateRequestReply(123, gomock.Any()).Return("", errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error creating reply: test error\n", func() {
			pssRequestReplyCreate(service, cmd, []string{"123"})
		})
	})

	t.Run("GetRequestError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestReplyCreateCmd()

		gomock.InOrder(
			service.EXPECT().CreateRequestReply(123, gomock.Any()).Return("C123", nil),
			service.EXPECT().GetReply("C123").Return(pss.Reply{}, errors.New("test error")),
		)

		test_output.AssertFatalOutput(t, "Error retrieving new reply: test error\n", func() {
			pssRequestReplyCreate(service, cmd, []string{"123"})
		})
	})
}
