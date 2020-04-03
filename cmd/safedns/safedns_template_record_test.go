package safedns

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func Test_safednsTemplateRecordListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := safednsTemplateRecordListCmd(nil)
		err := cmd.Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("MissingTemplate_Error", func(t *testing.T) {
		cmd := safednsTemplateRecordListCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_safednsTemplateRecordList(t *testing.T) {
	t.Run("Retrieve_ByTemplateID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetTemplateRecords(123, gomock.Any()).Return([]safedns.Record{}, nil).Times(1)

		safednsTemplateRecordList(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("Retrieve_ByTemplateName", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		expectedParameters := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				connection.APIRequestFiltering{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{"test template 1"},
				},
			},
		}

		gomock.InOrder(
			service.EXPECT().GetTemplates(gomock.Eq(expectedParameters)).Return([]safedns.Template{safedns.Template{ID: 123}}, nil),
			service.EXPECT().GetTemplateRecords(123, gomock.Any()).Return([]safedns.Record{}, nil).Times(1),
		)

		safednsTemplateRecordList(service, &cobra.Command{}, []string{"test template 1"})
	})

	t.Run("ExpectedFilterFromFlags", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateRecordListCmd(nil)
		cmd.Flags().Set("name", "test.testdomain1.co.uk")
		cmd.Flags().Set("type", "A")
		cmd.Flags().Set("content", "1.2.3.4")

		expectedParameters := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				connection.APIRequestFiltering{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{"test.testdomain1.co.uk"},
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

		service.EXPECT().GetTemplateRecords(123, gomock.Eq(expectedParameters)).Return([]safedns.Record{}, nil)

		safednsTemplateRecordList(service, cmd, []string{"123"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := safednsTemplateRecordList(service, cmd, []string{"123"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetTemplatesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetTemplates(gomock.Any()).Return([]safedns.Template{}, errors.New("test error")).Times(1)

		err := safednsTemplateRecordList(service, &cobra.Command{}, []string{"test template 1"})

		assert.Equal(t, "Error locating template [test template 1]: Error retrieving items: test error", err.Error())
	})

	t.Run("GetTemplateRecordsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetTemplateRecords(123, gomock.Any()).Return([]safedns.Record{}, errors.New("test error")).Times(1)

		err := safednsTemplateRecordList(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error retrieving records for template: test error", err.Error())
	})
}

func Test_safednsTemplateRecordShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := safednsTemplateRecordShowCmd(nil)
		err := cmd.Args(nil, []string{"123", "456"})

		assert.Nil(t, err)
	})

	t.Run("MissingTemplate_Error", func(t *testing.T) {
		cmd := safednsTemplateRecordShowCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})

	t.Run("MissingRecord_Error", func(t *testing.T) {
		cmd := safednsTemplateRecordShowCmd(nil)
		err := cmd.Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing record", err.Error())
	})
}

func Test_safednsTemplateRecordShow(t *testing.T) {
	t.Run("RetrieveSingle_ByTemplateID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetTemplateRecord(123, 456).Return(safedns.Record{}, nil).Times(1)

		safednsTemplateRecordShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("RetrieveMultiple_ByTemplateID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetTemplateRecord(123, 456).Return(safedns.Record{}, nil).Times(1),
			service.EXPECT().GetTemplateRecord(123, 789).Return(safedns.Record{}, nil).Times(1),
		)

		safednsTemplateRecordShow(service, &cobra.Command{}, []string{"123", "456", "789"})
	})

	t.Run("Retrieve_ByTemplateName", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		expectedParameters := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				connection.APIRequestFiltering{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{"test template 1"},
				},
			},
		}

		gomock.InOrder(
			service.EXPECT().GetTemplates(gomock.Eq(expectedParameters)).Return([]safedns.Template{safedns.Template{ID: 123}}, nil),
			service.EXPECT().GetTemplateRecord(123, 456).Return(safedns.Record{}, nil).Times(1),
		)

		safednsTemplateRecordShow(service, &cobra.Command{}, []string{"test template 1", "456"})
	})

	t.Run("GetTemplatesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetTemplates(gomock.Any()).Return([]safedns.Template{}, errors.New("test error")).Times(1)

		err := safednsTemplateRecordShow(service, &cobra.Command{}, []string{"test template 1"})

		assert.Equal(t, "Error locating template [test template 1]: Error retrieving items: test error", err.Error())
	})

	t.Run("InvalidRecordID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid record ID [abc]\n", func() {
			safednsTemplateRecordShow(service, &cobra.Command{}, []string{"123", "abc"})
		})
	})

	t.Run("GetTemplateRecordError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetTemplateRecord(123, 456).Return(safedns.Record{}, errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error retrieving record [456]: test error\n", func() {
			safednsTemplateRecordShow(service, &cobra.Command{}, []string{"123", "456"})
		})
	})
}

func Test_safednsTemplateRecordCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := safednsTemplateRecordCreateCmd(nil)
		err := cmd.Args(nil, []string{"123", "456"})

		assert.Nil(t, err)
	})

	t.Run("MissingTemplate_Error", func(t *testing.T) {
		cmd := safednsTemplateRecordCreateCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_safednsTemplateRecordCreate(t *testing.T) {
	t.Run("DefaultCreate_ByTemplateID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateRecordCreateCmd(nil)
		cmd.Flags().Set("name", "test.testdomain.co.uk")
		cmd.Flags().Set("type", "A")
		cmd.Flags().Set("content", "1.2.3.4")

		expectedRequest := safedns.CreateRecordRequest{
			Name:    "test.testdomain.co.uk",
			Type:    "A",
			Content: "1.2.3.4",
		}

		gomock.InOrder(
			service.EXPECT().CreateTemplateRecord(123, expectedRequest).Return(456, nil),
			service.EXPECT().GetTemplateRecord(123, 456).Return(safedns.Record{}, nil),
		)

		safednsTemplateRecordCreate(service, cmd, []string{"123"})
	})

	t.Run("DefaultCreate_ByTemplateName", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateRecordCreateCmd(nil)
		cmd.Flags().Set("name", "test.testdomain.co.uk")
		cmd.Flags().Set("type", "A")
		cmd.Flags().Set("content", "1.2.3.4")
		cmd.Flags().Set("priority", "0")

		expectedRequest := safedns.CreateRecordRequest{
			Name:     "test.testdomain.co.uk",
			Type:     "A",
			Content:  "1.2.3.4",
			Priority: ptr.Int(0),
		}

		gomock.InOrder(
			service.EXPECT().GetTemplates(gomock.Eq(
				connection.APIRequestParameters{
					Filtering: []connection.APIRequestFiltering{
						connection.APIRequestFiltering{
							Property: "name",
							Operator: connection.EQOperator,
							Value:    []string{"test template 1"},
						},
					},
				}),
			).Return([]safedns.Template{safedns.Template{ID: 123}}, nil),
			service.EXPECT().CreateTemplateRecord(123, expectedRequest).Return(456, nil),
			service.EXPECT().GetTemplateRecord(123, 456).Return(safedns.Record{}, nil),
		)

		safednsTemplateRecordCreate(service, cmd, []string{"test template 1"})
	})

	t.Run("GetTemplatesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetTemplates(gomock.Any()).Return([]safedns.Template{}, errors.New("test error")).Times(1)

		err := safednsTemplateRecordCreate(service, &cobra.Command{}, []string{"test template 1"})

		assert.Equal(t, "Error locating template [test template 1]: Error retrieving items: test error", err.Error())
	})

	t.Run("CreateTemplateRecordError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().CreateTemplateRecord(123, gomock.Any()).Return(456, errors.New("test error")),
		)

		err := safednsTemplateRecordCreate(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error creating record: test error", err.Error())
	})

	t.Run("GetTemplateRecordError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().CreateTemplateRecord(123, gomock.Any()).Return(456, nil),
			service.EXPECT().GetTemplateRecord(123, 456).Return(safedns.Record{}, errors.New("test error")),
		)

		err := safednsTemplateRecordCreate(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error retrieving new record: test error", err.Error())
	})
}

func Test_safednsTemplateRecordUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := safednsTemplateRecordUpdateCmd(nil)
		err := cmd.Args(nil, []string{"123", "456"})

		assert.Nil(t, err)
	})

	t.Run("MissingTemplate_Error", func(t *testing.T) {
		cmd := safednsTemplateRecordUpdateCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})

	t.Run("MissingRecord_Error", func(t *testing.T) {
		cmd := safednsTemplateRecordUpdateCmd(nil)
		err := cmd.Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing record", err.Error())
	})
}

func Test_safednsTemplateRecordUpdate(t *testing.T) {
	t.Run("UpdateSingle_ByTemplateID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchTemplateRecord(123, 456, gomock.Any()),
			service.EXPECT().GetTemplateRecord(123, 456).Return(safedns.Record{}, nil),
		)

		safednsTemplateRecordUpdate(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("UpdateMultiple_ByTemplateID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchTemplateRecord(123, 456, gomock.Any()).Return(456, nil),
			service.EXPECT().GetTemplateRecord(123, 456).Return(safedns.Record{}, nil),
			service.EXPECT().PatchTemplateRecord(123, 789, gomock.Any()).Return(789, nil),
			service.EXPECT().GetTemplateRecord(123, 789).Return(safedns.Record{}, nil),
		)

		safednsTemplateRecordUpdate(service, &cobra.Command{}, []string{"123", "456", "789"})
	})

	t.Run("ExpectedFlags", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateRecordUpdateCmd(nil)
		cmd.Flags().Set("name", "test.testdomain.co.uk")
		cmd.Flags().Set("type", "A")
		cmd.Flags().Set("content", "1.2.3.4")
		cmd.Flags().Set("priority", "0")

		expectedRequest := safedns.PatchRecordRequest{
			Name:     "test.testdomain.co.uk",
			Type:     "A",
			Content:  "1.2.3.4",
			Priority: ptr.Int(0),
		}

		gomock.InOrder(
			service.EXPECT().PatchTemplateRecord(123, 456, expectedRequest),
			service.EXPECT().GetTemplateRecord(123, 456).Return(safedns.Record{}, nil),
		)

		safednsTemplateRecordUpdate(service, cmd, []string{"123", "456"})
	})

	t.Run("UpdateSingle_ByTemplateName", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		gomock.InOrder(
			service.EXPECT().GetTemplates(gomock.Eq(connection.APIRequestParameters{
				Filtering: []connection.APIRequestFiltering{
					connection.APIRequestFiltering{
						Property: "name",
						Operator: connection.EQOperator,
						Value:    []string{"test template 1"},
					},
				},
			})).Return([]safedns.Template{safedns.Template{ID: 123}}, nil),
			service.EXPECT().PatchTemplateRecord(123, 456, gomock.Any()),
			service.EXPECT().GetTemplateRecord(123, 456).Return(safedns.Record{}, nil),
		)

		safednsTemplateRecordUpdate(service, &cobra.Command{}, []string{"test template 1", "456"})
	})

	t.Run("GetTemplatesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetTemplates(gomock.Any()).Return([]safedns.Template{}, errors.New("test error")).Times(1)

		err := safednsTemplateRecordUpdate(service, &cobra.Command{}, []string{"test template 1"})

		assert.Equal(t, "Error locating template [test template 1]: Error retrieving items: test error", err.Error())
	})

	t.Run("InvalidRecordID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid record ID [abc]\n", func() {
			safednsTemplateRecordUpdate(service, &cobra.Command{}, []string{"123", "abc"})
		})
	})

	t.Run("PatchTemplateRecordError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().PatchTemplateRecord(123, 456, gomock.Any()).Return(0, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating record [456]: test error\n", func() {
			safednsTemplateRecordUpdate(service, &cobra.Command{}, []string{"123", "456"})
		})
	})

	t.Run("GetTemplateRecordError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchTemplateRecord(123, 456, gomock.Any()),
			service.EXPECT().GetTemplateRecord(123, 456).Return(safedns.Record{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated record [456]: test error\n", func() {
			safednsTemplateRecordUpdate(service, &cobra.Command{}, []string{"123", "456"})
		})
	})
}

func Test_safednsTemplateRecordDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := safednsTemplateRecordDeleteCmd(nil)
		err := cmd.Args(nil, []string{"123", "456"})

		assert.Nil(t, err)
	})

	t.Run("MissingTemplate_Error", func(t *testing.T) {
		cmd := safednsTemplateRecordDeleteCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})

	t.Run("MissingRecord_Error", func(t *testing.T) {
		cmd := safednsTemplateRecordDeleteCmd(nil)
		err := cmd.Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing record", err.Error())
	})
}

func Test_safednsTemplateRecordDelete(t *testing.T) {
	t.Run("RetrieveSingle_ByTemplateID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().DeleteTemplateRecord(123, 456).Return(nil).Times(1)

		safednsTemplateRecordDelete(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("RetrieveMultiple_ByTemplateID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteTemplateRecord(123, 456).Return(nil).Times(1),
			service.EXPECT().DeleteTemplateRecord(123, 789).Return(nil).Times(1),
		)

		safednsTemplateRecordDelete(service, &cobra.Command{}, []string{"123", "456", "789"})
	})

	t.Run("Retrieve_ByTemplateName", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetTemplates(gomock.Eq(
				connection.APIRequestParameters{
					Filtering: []connection.APIRequestFiltering{
						connection.APIRequestFiltering{
							Property: "name",
							Operator: connection.EQOperator,
							Value:    []string{"test template 1"},
						},
					},
				}),
			).Return([]safedns.Template{safedns.Template{ID: 123}}, nil),
			service.EXPECT().DeleteTemplateRecord(123, 456).Return(nil).Times(1),
		)

		safednsTemplateRecordDelete(service, &cobra.Command{}, []string{"test template 1", "456"})
	})

	t.Run("GetTemplatesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetTemplates(gomock.Any()).Return([]safedns.Template{}, errors.New("test error")).Times(1)

		err := safednsTemplateRecordDelete(service, &cobra.Command{}, []string{"test template 1"})

		assert.Equal(t, "Error locating template [test template 1]: Error retrieving items: test error", err.Error())
	})
}
