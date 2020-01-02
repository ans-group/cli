package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
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

		err := pssRequestReplyList(service, &cobra.Command{}, []string{"abc"})
		assert.Equal(t, "Invalid request ID [abc]", err.Error())
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := pssRequestReplyList(service, cmd, []string{"123"})
		assert.Equal(t, "Missing value for filtering", err.Error())
	})

	t.Run("GetRequestConversationError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetRequestConversation(123, gomock.Any()).Return([]pss.Reply{}, errors.New("test error")).Times(1)

		err := pssRequestReplyList(service, &cobra.Command{}, []string{"123"})
		assert.Equal(t, "Error retrieving request replies: test error", err.Error())
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

		err := pssRequestReplyCreate(service, cmd, []string{"invalid"})
		assert.Contains(t, err.Error(), "Invalid request ID [invalid]")
	})

	t.Run("CreateRequestReplyError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestReplyCreateCmd()

		service.EXPECT().CreateRequestReply(123, gomock.Any()).Return("", errors.New("test error")).Times(1)

		err := pssRequestReplyCreate(service, cmd, []string{"123"})
		assert.Equal(t, "Error creating reply: test error", err.Error())
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

		err := pssRequestReplyCreate(service, cmd, []string{"123"})
		assert.Equal(t, "Error retrieving new reply: test error", err.Error())
	})
}
