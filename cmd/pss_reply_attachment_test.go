package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/test"
)

func Test_pssReplyAttachmentDownloadCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssReplyAttachmentDownloadCmd().Args(nil, []string{"123", "test.txt"})

		assert.Nil(t, err)
	})

	t.Run("MissingReply_Error", func(t *testing.T) {
		err := pssReplyAttachmentDownloadCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing reply", err.Error())
	})

	t.Run("MissingAttachment_Error", func(t *testing.T) {
		err := pssReplyAttachmentDownloadCmd().Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing attachment", err.Error())
	})
}

func Test_pssReplyAttachmentDownload(t *testing.T) {
	t.Run("Valid_DownloadsFile", func(t *testing.T) {
		appFilesystem = afero.NewMemMapFs()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		attachmentStream := ioutil.NopCloser(bytes.NewReader([]byte("test content")))

		service := mocks.NewMockPSSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DownloadReplyAttachmentStream("C123", "test1.txt").Return(attachmentStream, nil),
		)

		pssReplyAttachmentDownload(service, pssReplyAttachmentDownloadCmd(), []string{"C123", "test1.txt"})
	})

	t.Run("DownloadReplyAttachmentStreamError_ReturnsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DownloadReplyAttachmentStream("C123", "test1.txt").Return(nil, errors.New("test error")),
		)

		test_output.AssertFatalOutputFunc(t, func(stdErr string) {
			assert.Equal(t, stdErr, "Error downloading reply attachment: test error\n")
		}, func() {
			pssReplyAttachmentDownload(service, pssReplyAttachmentDownloadCmd(), []string{"C123", "test1.txt"})
		})
	})

	t.Run("FileExists_ReturnsFatal", func(t *testing.T) {
		appFilesystem = afero.NewMemMapFs()
		afero.WriteFile(appFilesystem, "test1.txt", []byte{}, 0644)

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DownloadReplyAttachmentStream("C123", "test1.txt").Return(nil, nil),
		)

		test_output.AssertFatalOutputFunc(t, func(stdErr string) {
			assert.Equal(t, stdErr, "Destination file [test1.txt] exists\n")
		}, func() {
			pssReplyAttachmentDownload(service, pssReplyAttachmentDownloadCmd(), []string{"C123", "test1.txt"})
		})
	})

	t.Run("WriteReaderError_ReturnsFatal", func(t *testing.T) {
		appFilesystem = afero.NewMemMapFs()
		b := test.TestReadCloser{
			ReadError: errors.New("test reader error 1"),
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssReplyAttachmentDownloadCmd()
		cmd.Flags().Set("path", "/some/path/test.txt")

		gomock.InOrder(
			service.EXPECT().DownloadReplyAttachmentStream("C123", "test1.txt").Return(&b, nil),
		)

		test_output.AssertFatalOutputFunc(t, func(stdErr string) {
			assert.Contains(t, stdErr, "test reader error 1")
		}, func() {
			pssReplyAttachmentDownload(service, cmd, []string{"C123", "test1.txt"})
		})
	})
}
