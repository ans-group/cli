package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
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
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		test_output.AssertFatalOutput(t, "Missing value for filtering\n", func() {
			ddosxDomainRecordList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainRecordsError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainRecords("testdomain1.co.uk", gomock.Any()).Return([]ddosx.Record{}, errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving domain records: test error\n", func() {
			ddosxDomainRecordList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainRecordShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainRecordShowCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainRecordShowCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRecord_Error", func(t *testing.T) {
		cmd := ddosxDomainRecordShowCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing record", err.Error())
	})
}

func Test_ddosxDomainRecordShow(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.Record{}, nil)

		ddosxDomainRecordShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("GetDomainRecord_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.Record{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving domain record [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainRecordShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
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

		gomock.InOrder(
			service.EXPECT().CreateDomainRecord("testdomain1.co.uk", gomock.Eq(expectedRequest)).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.Record{}, nil),
		)

		ddosxDomainRecordCreate(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("CreateDomainRecordError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().CreateDomainRecord("testdomain1.co.uk", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error creating domain record: test error\n", func() {
			ddosxDomainRecordCreate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainRecordError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().CreateDomainRecord("testdomain1.co.uk", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.Record{}, errors.New("test error")),
		)

		test_output.AssertFatalOutput(t, "Error retrieving new domain record [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainRecordCreate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
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

		gomock.InOrder(
			service.EXPECT().PatchDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.Record{}, nil),
		)

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

		gomock.InOrder(
			service.EXPECT().PatchDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.Record{}, nil),
		)

		ddosxDomainRecordUpdate(service, cmd, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("UpdateDomainRecordError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().PatchDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating domain record [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainRecordUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})

	t.Run("GetDomainRecordError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(nil),
			service.EXPECT().GetDomainRecord("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.Record{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated domain record [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainRecordUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
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

		test_output.AssertErrorOutput(t, "Error removing domain record [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainRecordDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}
