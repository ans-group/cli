package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func Test_ddosxDomainRecordListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainRecordListCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainRecordListCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainRecordList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainRecords("testdomain1.co.uk", gomock.Any()).Return([]ddosx.Record{}, nil).Times(1)

		ddosxDomainRecordList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		output := test.CatchStdErr(t, func() {
			ddosxDomainRecordList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Missing value for filtering\n", output)
	})

	t.Run("GetDomainRecordsError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainRecords("testdomain1.co.uk", gomock.Any()).Return([]ddosx.Record{}, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			ddosxDomainRecordList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving domain records: test error\n", output)
	})
}

func Test_ddosxDomainRecordCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainRecordCreateCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainRecordCreateCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainRecordCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainRecordCreateCmd()
		cmd.Flags().Set("name", "sub.testdomain1.co.uk")

		expectedRequest := ddosx.CreateRecordRequest{
			Name: "sub.testdomain1.co.uk",
		}

		service.EXPECT().CreateDomainRecord("testdomain1.co.uk", gomock.Eq(expectedRequest)).Return("00000000-0000-0000-0000-000000000000", nil).Times(1)

		ddosxDomainRecordCreate(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("CreateDomainRecord_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().CreateDomainRecord("testdomain1.co.uk", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			ddosxDomainRecordCreate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error creating domain record: test error\n", output)
	})
}

func Test_ddosxDomainRecordUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainRecordUpdateCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainRecordUpdateCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRecord_Error", func(t *testing.T) {
		cmd := ddosxDomainRecordUpdateCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing record", err.Error())
	})
}

func Test_ddosxDomainRecordUpdate(t *testing.T) {
	t.Run("DefaultUpdate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainRecordUpdateCmd()
		cmd.Flags().Set("content", "1.2.3.4")

		expectedRequest := ddosx.PatchRecordRequest{
			Content: "1.2.3.4",
		}

		service.EXPECT().PatchDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return(nil)

		ddosxDomainRecordUpdate(service, cmd, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("AllFlagsSet", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainRecordUpdateCmd()
		cmd.Flags().Set("name", "sub1.testdomain1.co.uk")
		cmd.Flags().Set("type", "A")
		cmd.Flags().Set("content", "1.2.3.4")
		cmd.Flags().Set("ssl-id", "00000000-0000-0000-0000-000000000000")
		cmd.Flags().Set("safedns-record-id", "123")

		expectedRequest := ddosx.PatchRecordRequest{
			Name:            "sub1.testdomain1.co.uk",
			Type:            ddosx.RecordTypeA,
			Content:         "1.2.3.4",
			SSLID:           "00000000-0000-0000-0000-000000000000",
			SafeDNSRecordID: 123,
		}

		service.EXPECT().PatchDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return(nil)

		ddosxDomainRecordUpdate(service, cmd, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("UpdateDomainRecord_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().PatchDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			ddosxDomainRecordUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})

		assert.Equal(t, "Error updating domain record [00000000-0000-0000-0000-000000000000]: test error\n", output)
	})
}

func Test_ddosxDomainRecordDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainRecordDeleteCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainRecordDeleteCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRecord_Error", func(t *testing.T) {
		cmd := ddosxDomainRecordDeleteCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing record", err.Error())
	})
}

func Test_ddosxDomainRecordDelete(t *testing.T) {
	t.Run("DefaultDelete", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(nil)

		ddosxDomainRecordDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("DeleteDomainRecordError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			ddosxDomainRecordDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})

		assert.Equal(t, "Error removing domain record [00000000-0000-0000-0000-000000000000]: test error\n", output)
	})
}
