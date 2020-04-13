package pss

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/test"
)

func Test_pssReplyAttachmentDownloadCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssReplyAttachmentDownloadCmd(nil, nil).Args(nil, []string{"123", "test.txt"})

		assert.Nil(t, err)
	})

	t.Run("MissingReply_Error", func(t *testing.T) {
		err := pssReplyAttachmentDownloadCmd(nil, nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing reply", err.Error())
	})

	t.Run("MissingAttachment_Error", func(t *testing.T) {
		err := pssReplyAttachmentDownloadCmd(nil, nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing attachment", err.Error())
	})
}

func Test_pssReplyAttachmentDownload(t *testing.T) {
	t.Run("Valid_DownloadsFile", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		attachmentStream := ioutil.NopCloser(bytes.NewReader([]byte("test content")))

		service := mocks.NewMockPSSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DownloadReplyAttachmentStream("C123", "test1.txt").Return(attachmentStream, nil),
		)

		pssReplyAttachmentDownload(service, fs, &cobra.Command{}, []string{"C123", "test1.txt"})
	})

	t.Run("DownloadReplyAttachmentStreamError_ReturnsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DownloadReplyAttachmentStream("C123", "test1.txt").Return(nil, errors.New("test error")),
		)

		err := pssReplyAttachmentDownload(service, nil, &cobra.Command{}, []string{"C123", "test1.txt"})
		assert.Equal(t, "Error downloading reply attachment: test error", err.Error())
	})

	t.Run("FileExists_ReturnsFatal", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		afero.WriteFile(fs, "test1.txt", []byte{}, 0644)

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DownloadReplyAttachmentStream("C123", "test1.txt").Return(nil, nil),
		)

		err := pssReplyAttachmentDownload(service, fs, &cobra.Command{}, []string{"C123", "test1.txt"})
		assert.Equal(t, "Destination file [test1.txt] exists", err.Error())
	})

	t.Run("WriteReaderError_ReturnsFatal", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		b := test.TestReadCloser{
			ReadError: errors.New("test reader error 1"),
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := pssReplyAttachmentDownloadCmd(nil, nil)
		cmd.Flags().Set("path", "/some/path/test.txt")

		gomock.InOrder(
			service.EXPECT().DownloadReplyAttachmentStream("C123", "test1.txt").Return(&b, nil),
		)

		err := pssReplyAttachmentDownload(service, fs, cmd, []string{"C123", "test1.txt"})
		assert.Contains(t, err.Error(), "test reader error 1")
	})
}

func Test_pssReplyAttachmentUploadCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssReplyAttachmentUploadCmd(nil, nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("MissingReply_Error", func(t *testing.T) {
		err := pssReplyAttachmentUploadCmd(nil, nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing reply", err.Error())
	})
}

func Test_pssReplyAttachmentUpload(t *testing.T) {
	t.Run("Valid_UploadsFile", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		afero.WriteFile(fs, "/test/test1.txt", []byte("test content"), 644)

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		cmd := pssReplyAttachmentUploadCmd(nil, nil)
		cmd.Flags().Set("path", "/test/test1.txt")

		gomock.InOrder(
			service.EXPECT().UploadReplyAttachmentStream("C123", "test1.txt", gomock.Any()).Return(nil),
		)

		pssReplyAttachmentUpload(service, fs, cmd, []string{"C123"})
	})

	t.Run("FileOpenError_ReturnsError", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		cmd := pssReplyAttachmentUploadCmd(nil, nil)
		cmd.Flags().Set("path", "/test/test1.txt")

		err := pssReplyAttachmentUpload(nil, fs, cmd, []string{"C123"})

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "Failed to open file")
	})

	t.Run("UploadReplyAttachmentStream_ReturnsError", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		afero.WriteFile(fs, "/test/test1.txt", []byte("test content"), 644)

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		cmd := pssReplyAttachmentUploadCmd(nil, nil)
		cmd.Flags().Set("path", "/test/test1.txt")

		gomock.InOrder(
			service.EXPECT().UploadReplyAttachmentStream("C123", "test1.txt", gomock.Any()).Return(errors.New("test error")),
		)

		err := pssReplyAttachmentUpload(service, fs, cmd, []string{"C123"})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "Failed to upload attachment")
	})
}

func Test_pssReplyAttachmentDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssReplyAttachmentDeleteCmd(nil).Args(nil, []string{"123", "test.txt"})

		assert.Nil(t, err)
	})

	t.Run("MissingReply_Error", func(t *testing.T) {
		err := pssReplyAttachmentDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing reply", err.Error())
	})

	t.Run("MissingAttachment_Error", func(t *testing.T) {
		err := pssReplyAttachmentDeleteCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing attachment", err.Error())
	})
}

func Test_pssReplyAttachmentDelete(t *testing.T) {
	t.Run("Valid_DeletesFile", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteReplyAttachment("C123", "test1.txt").Return(nil),
		)

		pssReplyAttachmentDelete(service, pssReplyAttachmentDeleteCmd(nil), []string{"C123", "test1.txt"})
	})

	t.Run("DeleteReplyAttachmentError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteReplyAttachment("C123", "test1.txt").Return(errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error deleting reply attachment [test1.txt]: test error\n", func() {
			pssReplyAttachmentDelete(service, pssReplyAttachmentDeleteCmd(nil), []string{"C123", "test1.txt"})
		})
	})
}
