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
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func Test_safednsZoneRecordListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := safednsZoneRecordListCmd()
		err := cmd.Args(nil, []string{"testdomain1.com"})

		assert.Nil(t, err)
	})

	t.Run("MissingZone_Error", func(t *testing.T) {
		cmd := safednsZoneRecordListCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})
}

func Test_safednsZoneRecordList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetZoneRecords("testdomain1.com", gomock.Any()).Return([]safedns.Record{}, nil).Times(1)

		safednsZoneRecordList(service, &cobra.Command{}, []string{"testdomain1.com"})
	})

	t.Run("ExpectedFilterFromFlags", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsZoneRecordListCmd()
		cmd.Flags().Set("name", "testdomain1.com")
		cmd.Flags().Set("type", "A")
		cmd.Flags().Set("content", "1.2.3.4")

		expectedParams := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				connection.APIRequestFiltering{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{"testdomain1.com"},
				},
				connection.APIRequestFiltering{
					Property: "type",
					Operator: connection.EQOperator,
					Value:    []string{"A"},
				},
				connection.APIRequestFiltering{
					Property: "content",
					Operator: connection.EQOperator,
					Value:    []string{"1.2.3.4"},
				},
			},
		}

		service.EXPECT().GetZoneRecords("testdomain1.com", gomock.Eq(expectedParams)).Return([]safedns.Record{}, nil)

		safednsZoneRecordList(service, cmd, []string{"testdomain1.com"})
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

		service := mocks.NewMockSafeDNSService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		output := test.CatchStdErr(t, func() {
			safednsZoneRecordList(service, &cobra.Command{}, []string{"123"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Missing value for filtering\n", output)
	})

	t.Run("GetZonesError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetZoneRecords("testdomain1.com", gomock.Any()).Return([]safedns.Record{}, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			safednsZoneRecordList(service, &cobra.Command{}, []string{"testdomain1.com"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving records for zone: test error\n", output)
	})
}

func Test_safednsZoneRecordShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := safednsZoneRecordShowCmd()
		err := cmd.Args(nil, []string{"testdomain1.com", "123"})

		assert.Nil(t, err)
	})

	t.Run("MissingZone_Error", func(t *testing.T) {
		cmd := safednsZoneRecordShowCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})

	t.Run("MissingRecord_Error", func(t *testing.T) {
		cmd := safednsZoneRecordShowCmd()
		err := cmd.Args(nil, []string{"testdomain1.com"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing record", err.Error())
	})
}

func Test_safednsZoneRecordShow(t *testing.T) {
	t.Run("SingleZoneRecord", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetZoneRecord("testdomain1.com", 123).Return(safedns.Record{}, nil).Times(1)

		safednsZoneRecordShow(service, &cobra.Command{}, []string{"testdomain1.com", "123"})
	})

	t.Run("MultipleZoneRecords", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetZoneRecord("testdomain1.com", 123).Return(safedns.Record{}, nil),
			service.EXPECT().GetZoneRecord("testdomain1.com", 456).Return(safedns.Record{}, nil),
		)

		safednsZoneRecordShow(service, &cobra.Command{}, []string{"testdomain1.com", "123", "456"})
	})

	t.Run("InvalidRecordID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		output := test.CatchStdErr(t, func() {
			safednsZoneRecordShow(service, &cobra.Command{}, []string{"testdomain1.com", "abc"})
		})

		assert.Equal(t, "Invalid record ID [abc]\n", output)
	})

	t.Run("GetZoneRecordError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetZoneRecord("testdomain1.com", 123).Return(safedns.Record{}, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			safednsZoneRecordShow(service, &cobra.Command{}, []string{"testdomain1.com", "123"})
		})

		assert.Equal(t, "Error retrieving record [123]: test error\n", output)
	})
}

func Test_safednsZoneRecordCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := safednsZoneRecordCreateCmd()
		err := cmd.Args(nil, []string{"testdomain1.com"})

		assert.Nil(t, err)
	})

	t.Run("MissingZone_Error", func(t *testing.T) {
		cmd := safednsZoneRecordCreateCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})
}

func Test_safednsZoneRecordCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsZoneRecordCreateCmd()
		cmd.Flags().Set("name", "www.testdomain1.com")
		cmd.Flags().Set("type", "A")
		cmd.Flags().Set("content", "1.2.3.4")

		expectedRequest := safedns.CreateRecordRequest{
			Name:    "www.testdomain1.com",
			Type:    "A",
			Content: "1.2.3.4",
		}

		gomock.InOrder(
			service.EXPECT().CreateZoneRecord("testdomain1.com", expectedRequest).Return(123, nil),
			service.EXPECT().GetZoneRecord("testdomain1.com", 123).Return(safedns.Record{}, nil),
		)

		safednsZoneRecordCreate(service, cmd, []string{"testdomain1.com"})
	})

	t.Run("DefaultCreate_WithPriorityDefaultValue", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsZoneRecordCreateCmd()
		cmd.Flags().Set("priority", "0")

		gomock.InOrder(
			service.EXPECT().CreateZoneRecord("testdomain1.com", gomock.Any()).Return(123, nil).Do(func(zoneName string, req safedns.CreateRecordRequest) {
				if req.Priority == nil {
					t.Fatal("Expected non-nil priority")
				}
				if *req.Priority != 0 {
					t.Fatalf("Expected priority of 0, got %d", *req.Priority)
				}
			}),
			service.EXPECT().GetZoneRecord("testdomain1.com", 123).Return(safedns.Record{}, nil),
		)

		safednsZoneRecordCreate(service, cmd, []string{"testdomain1.com"})
	})

	t.Run("CreateZoneRecordError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().CreateZoneRecord("testdomain1.com", gomock.Any()).Return(0, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			safednsZoneRecordCreate(service, &cobra.Command{}, []string{"testdomain1.com"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error creating record: test error\n", output)
	})

	t.Run("GetZoneRecordError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})

		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().CreateZoneRecord("testdomain1.com", gomock.Any()).Return(123, nil),
			service.EXPECT().GetZoneRecord("testdomain1.com", 123).Return(safedns.Record{}, errors.New("test error")),
		)

		output := test.CatchStdErr(t, func() {
			safednsZoneRecordCreate(service, &cobra.Command{}, []string{"testdomain1.com"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving new record: test error\n", output)
	})
}

func Test_safednsZoneRecordUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := safednsZoneRecordUpdateCmd()
		err := cmd.Args(nil, []string{"testdomain1.com", "123"})

		assert.Nil(t, err)
	})

	t.Run("MissingZone_Error", func(t *testing.T) {
		cmd := safednsZoneRecordUpdateCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})

	t.Run("MissingRecord_Error", func(t *testing.T) {
		cmd := safednsZoneRecordUpdateCmd()
		err := cmd.Args(nil, []string{"testdomain1.com"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing record", err.Error())
	})
}

func Test_safednsZoneRecordUpdate(t *testing.T) {
	t.Run("UpdateSingle", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsZoneRecordUpdateCmd()
		cmd.Flags().Set("name", "www.testdomain1.com")
		cmd.Flags().Set("type", "A")
		cmd.Flags().Set("content", "1.2.3.4")

		expectedRequest := safedns.PatchRecordRequest{
			Name:    "www.testdomain1.com",
			Type:    "A",
			Content: "1.2.3.4",
		}

		gomock.InOrder(
			service.EXPECT().PatchZoneRecord("testdomain1.com", 123, expectedRequest).Return(123, nil),
			service.EXPECT().GetZoneRecord("testdomain1.com", 123).Return(safedns.Record{}, nil),
		)

		safednsZoneRecordUpdate(service, cmd, []string{"testdomain1.com", "123"})
	})

	t.Run("UpdateMultiple", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsZoneRecordUpdateCmd()
		cmd.Flags().Set("name", "www.testdomain1.com")
		cmd.Flags().Set("type", "A")
		cmd.Flags().Set("content", "1.2.3.4")

		expectedRequest := safedns.PatchRecordRequest{
			Name:    "www.testdomain1.com",
			Type:    "A",
			Content: "1.2.3.4",
		}

		gomock.InOrder(
			service.EXPECT().PatchZoneRecord("testdomain1.com", 123, expectedRequest).Return(123, nil),
			service.EXPECT().GetZoneRecord("testdomain1.com", 123).Return(safedns.Record{}, nil),
			service.EXPECT().PatchZoneRecord("testdomain1.com", 456, expectedRequest).Return(456, nil),
			service.EXPECT().GetZoneRecord("testdomain1.com", 456).Return(safedns.Record{}, nil),
		)

		safednsZoneRecordUpdate(service, cmd, []string{"testdomain1.com", "123", "456"})
	})

	t.Run("Update_WithPriorityDefaultValue", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsZoneRecordUpdateCmd()
		cmd.Flags().Set("priority", "0")

		gomock.InOrder(
			service.EXPECT().PatchZoneRecord("testdomain1.com", 123, gomock.Any()).Return(123, nil).Do(func(zoneName string, recordID int, req safedns.PatchRecordRequest) {
				if req.Priority == nil || *req.Priority != 0 {
					t.Fatal("Unexpected record priority")
				}
			}),
			service.EXPECT().GetZoneRecord("testdomain1.com", 123).Return(safedns.Record{}, nil),
		)

		safednsZoneRecordUpdate(service, cmd, []string{"testdomain1.com", "123"})
	})

	t.Run("InvalidRecordID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		output := test.CatchStdErr(t, func() {
			safednsZoneRecordUpdate(service, &cobra.Command{}, []string{"testdomain1.com", "abc"})
		})

		assert.Equal(t, "Invalid record ID [abc]\n", output)
	})

	t.Run("PatchZoneRecordError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().PatchZoneRecord("testdomain1.com", 123, gomock.Any()).Return(0, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			safednsZoneRecordUpdate(service, &cobra.Command{}, []string{"testdomain1.com", "123"})
		})

		assert.Equal(t, "Error updating record [123]: test error\n", output)
	})

	t.Run("GetZoneRecordError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchZoneRecord("testdomain1.com", 123, gomock.Any()).Return(123, nil),
			service.EXPECT().GetZoneRecord("testdomain1.com", 123).Return(safedns.Record{}, errors.New("test error")),
		)

		output := test.CatchStdErr(t, func() {
			safednsZoneRecordUpdate(service, &cobra.Command{}, []string{"testdomain1.com", "123"})
		})

		assert.Equal(t, "Error retrieving updated record [123]: test error\n", output)
	})
}

func Test_safednsZoneRecordDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := safednsZoneRecordDeleteCmd()
		err := cmd.Args(nil, []string{"testdomain1.com", "123"})

		assert.Nil(t, err)
	})

	t.Run("MissingZone_Error", func(t *testing.T) {
		cmd := safednsZoneRecordDeleteCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})

	t.Run("MissingRecord_Error", func(t *testing.T) {
		cmd := safednsZoneRecordDeleteCmd()
		err := cmd.Args(nil, []string{"testdomain1.com"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing record", err.Error())
	})
}

func Test_safednsZoneRecordDelete(t *testing.T) {
	t.Run("SingleZoneRecord", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().DeleteZoneRecord("testdomain1.com", 123).Return(nil).Times(1)

		safednsZoneRecordDelete(service, &cobra.Command{}, []string{"testdomain1.com", "123"})
	})

	t.Run("MultipleZoneRecords", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteZoneRecord("testdomain1.com", 123).Return(nil),
			service.EXPECT().DeleteZoneRecord("testdomain1.com", 456).Return(nil),
		)

		safednsZoneRecordDelete(service, &cobra.Command{}, []string{"testdomain1.com", "123", "456"})
	})

	t.Run("InvalidRecordID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		output := test.CatchStdErr(t, func() {
			safednsZoneRecordDelete(service, &cobra.Command{}, []string{"testdomain1.com", "abc"})
		})

		assert.Equal(t, "Invalid record ID [abc]\n", output)
	})

	t.Run("DeleteZoneRecordError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().DeleteZoneRecord("testdomain1.com", 123).Return(errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			safednsZoneRecordDelete(service, &cobra.Command{}, []string{"testdomain1.com", "123"})
		})

		assert.Equal(t, "Error removing record [123]: test error\n", output)
	})
}
