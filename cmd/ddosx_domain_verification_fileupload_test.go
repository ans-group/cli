package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"

	"github.com/spf13/afero"
)

func Test_ddosxDomainVerificationFileUploadShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainVerificationFileUploadShowCmd().Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainVerificationFileUploadShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainVerificationFileUploadShow(t *testing.T) {
	t.Run("SingleDomainVerificationFileUpload", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DownloadDomainVerificationFile("testdomain1.co.uk").Return("testfilecontent", "testfilename", nil)

		ddosxDomainVerificationFileUploadShow(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MultipleDomainVerificationFileUploads", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DownloadDomainVerificationFile("testdomain1.co.uk").Return("testfilecontent", "testfilename", nil),
			service.EXPECT().DownloadDomainVerificationFile("testdomain2.co.uk").Return("testfilecontent", "testfilename", nil),
		)

		ddosxDomainVerificationFileUploadShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "testdomain2.co.uk"})
	})

	t.Run("DownloadDomainVerificationFileError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DownloadDomainVerificationFile("testdomain1.co.uk").Return("testfilecontent", "testfilename", errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving domain verification file [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainVerificationFileUploadShow(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainVerificationFileUploadDownloadCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainVerificationFileUploadDownloadCmd().Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainVerificationFileUploadDownloadCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainVerificationFileUploadDownload(t *testing.T) {
	t.Run("Valid_FileCreated", func(t *testing.T) {
		appFilesystem = afero.NewMemMapFs()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainVerificationFileUploadDownloadCmd()
		cmd.Flags().Set("path", "/tmp")

		service.EXPECT().DownloadDomainVerificationFile("testdomain1.co.uk").Return("testfilecontent", "testfilename.txt", nil)

		output := test.CatchStdOut(t, func() {
			ddosxDomainVerificationFileUploadDownload(service, cmd, []string{"testdomain1.co.uk"})
		})

		filename := filepath.Join("/tmp", "testfilename.txt")

		_, err := appFilesystem.Stat(filename)
		if os.IsNotExist(err) {
			t.Errorf("file \"%s\" does not exist.\n", filename)
		}

		assert.Equal(t, filename+"\n", output)
	})

	t.Run("FileExists_OutputsFatal", func(t *testing.T) {

		appFilesystem = afero.NewMemMapFs()
		afero.WriteFile(appFilesystem, "/tmp/testfilename.txt", []byte{}, 0644)

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainVerificationFileUploadDownloadCmd()
		cmd.Flags().Set("path", "/tmp")

		service.EXPECT().DownloadDomainVerificationFile("testdomain1.co.uk").Return("testfilecontent", "testfilename.txt", nil)

		filename := filepath.Join("/tmp", "testfilename.txt")
		test_output.AssertFatalOutput(t, fmt.Sprintf("Destination file [%s] exists\n", filename), func() {
			ddosxDomainVerificationFileUploadDownload(service, cmd, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("WriteFileError_OutputsFatal", func(t *testing.T) {

		appFilesystem = afero.NewRegexpFs(afero.NewMemMapFs(), regexp.MustCompile(`\.invalid$`))

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainVerificationFileUploadDownloadCmd()
		cmd.Flags().Set("path", "/tmp")

		service.EXPECT().DownloadDomainVerificationFile("testdomain1.co.uk").Return("testfilecontent", "testfilename.txt", nil)

		filename := filepath.Join("/tmp", "testfilename.txt")
		test_output.AssertFatalOutput(t, fmt.Sprintf("Error writing domain verification file to [%s]: open %s: file does not exist\n", filename, filename), func() {
			ddosxDomainVerificationFileUploadDownload(service, cmd, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("DownloadDomainVerificationFileError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DownloadDomainVerificationFile("testdomain1.co.uk").Return("testfilecontent", "testfilename", errors.New("test error"))

		test_output.AssertFatalOutput(t, "Error retrieving domain verification file: test error\n", func() {
			ddosxDomainVerificationFileUploadDownload(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainVerificationFileUploadVerifyCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainVerificationFileUploadVerifyCmd().Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainVerificationFileUploadVerifyCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainVerificationFileUploadVerify(t *testing.T) {
	t.Run("SingleDomain_NoError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().VerifyDomainFileUpload("testdomain1.co.uk").Return(nil)

		ddosxDomainVerificationFileUploadVerify(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("DownloadDomainVerificationFileError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().VerifyDomainFileUpload("testdomain1.co.uk").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error verifying domain [testdomain1.co.uk] via verification file method: test error\n", func() {
			ddosxDomainVerificationFileUploadVerify(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}
