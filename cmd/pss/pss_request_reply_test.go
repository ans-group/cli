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

func Test_pssRequestReplyListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssRequestReplyListCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := pssRequestReplyListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing request", err.Error())
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

	t.Run("GetRequestID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		err := pssRequestReplyList(service, &cobra.Command{}, []string{"abc"})
		assert.Equal(t, "invalid request ID [abc]", err.Error())
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := pssRequestReplyList(service, cmd, []string{"123"})
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetRequestConversationError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetRequestConversation(123, gomock.Any()).Return([]pss.Reply{}, errors.New("test error")).Times(1)

		err := pssRequestReplyList(service, &cobra.Command{}, []string{"123"})
		assert.Equal(t, "error retrieving request replies: test error", err.Error())
	})
}

func Test_pssRequestReplyCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssRequestReplyCreateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := pssRequestReplyCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing request", err.Error())
	})
}

func Test_pssRequestReplyCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestReplyCreateCmd(nil)
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

	t.Run("InvalidRequestID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestReplyCreateCmd(nil)

		err := pssRequestReplyCreate(service, cmd, []string{"invalid"})
		assert.Contains(t, err.Error(), "invalid request ID [invalid]")
	})

	t.Run("CreateRequestReplyError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestReplyCreateCmd(nil)

		service.EXPECT().CreateRequestReply(123, gomock.Any()).Return("", errors.New("test error")).Times(1)

		err := pssRequestReplyCreate(service, cmd, []string{"123"})
		assert.Equal(t, "error creating reply: test error", err.Error())
	})

	t.Run("GetRequestError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssRequestReplyCreateCmd(nil)

		gomock.InOrder(
			service.EXPECT().CreateRequestReply(123, gomock.Any()).Return("C123", nil),
			service.EXPECT().GetReply("C123").Return(pss.Reply{}, errors.New("test error")),
		)

		err := pssRequestReplyCreate(service, cmd, []string{"123"})
		assert.Equal(t, "error retrieving new reply: test error", err.Error())
	})
}
