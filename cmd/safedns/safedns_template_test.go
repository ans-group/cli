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
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func Test_safednsTemplateList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetTemplates(gomock.Any()).Return([]safedns.Template{}, nil).Times(1)

		safednsTemplateList(service, &cobra.Command{}, []string{})
	})

	t.Run("ExpectedFilterFromFlags", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateListCmd(nil)
		cmd.Flags().Set("name", "test template 1")

		expectedParams := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				connection.APIRequestFiltering{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{"test template 1"},
				},
			},
		}

		service.EXPECT().GetTemplates(gomock.Eq(expectedParams)).Return([]safedns.Template{}, nil).Times(1)

		safednsTemplateList(service, cmd, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := safednsTemplateList(service, cmd, []string{"123"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetTemplatesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetTemplates(gomock.Any()).Return([]safedns.Template{}, errors.New("test error")).Times(1)

		err := safednsTemplateList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving templates: test error", err.Error())
	})
}

func Test_safednsTemplateShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := safednsTemplateShowCmd(nil).Args(nil, []string{"test template 1"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := safednsTemplateShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_safednsTemplateShow(t *testing.T) {
	t.Run("SingleTemplate_ByID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetTemplate(123).Return(safedns.Template{}, nil).Times(1)

		safednsTemplateShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleTemplates_ByID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetTemplate(123).Return(safedns.Template{}, nil),
			service.EXPECT().GetTemplate(456).Return(safedns.Template{}, nil),
		)

		safednsTemplateShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("SingleTemplate_ByName", func(t *testing.T) {
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
		)

		safednsTemplateShow(service, &cobra.Command{}, []string{"test template 1"})
	})

	t.Run("MultipleTemplates_ByName", func(t *testing.T) {
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
			service.EXPECT().GetTemplates(gomock.Eq(
				connection.APIRequestParameters{
					Filtering: []connection.APIRequestFiltering{
						connection.APIRequestFiltering{
							Property: "name",
							Operator: connection.EQOperator,
							Value:    []string{"test template 2"},
						},
					},
				}),
			).Return([]safedns.Template{safedns.Template{ID: 456}}, nil),
		)

		safednsTemplateShow(service, &cobra.Command{}, []string{"test template 1", "test template 2"})
	})

	t.Run("GetTemplateError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetTemplate(123).Return(safedns.Template{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving template [123]: Error retrieving template by ID [123]: test error\n", func() {
			safednsTemplateShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_safednsTemplateCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateCreateCmd(nil)
		cmd.Flags().Set("name", "test template 1")
		cmd.Flags().Set("default", "true")

		expectedRequest := safedns.CreateTemplateRequest{
			Name:    "test template 1",
			Default: true,
		}

		gomock.InOrder(
			service.EXPECT().CreateTemplate(expectedRequest).Return(123, nil),
			service.EXPECT().GetTemplate(123).Return(safedns.Template{}, nil),
		)

		safednsTemplateCreate(service, cmd, []string{})
	})

	t.Run("CreateTemplateError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateCreateCmd(nil)
		cmd.Flags().Set("name", "test template 1")
		cmd.Flags().Set("default", "true")

		expectedRequest := safedns.CreateTemplateRequest{
			Name:    "test template 1",
			Default: true,
		}

		service.EXPECT().CreateTemplate(expectedRequest).Return(0, errors.New("test error")).Times(1)

		err := safednsTemplateCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating template: test error", err.Error())
	})

	t.Run("GetTemplateError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateCreateCmd(nil)
		cmd.Flags().Set("name", "test template 1")
		cmd.Flags().Set("default", "true")

		expectedRequest := safedns.CreateTemplateRequest{
			Name:    "test template 1",
			Default: true,
		}

		gomock.InOrder(
			service.EXPECT().CreateTemplate(expectedRequest).Return(123, nil),
			service.EXPECT().GetTemplate(123).Return(safedns.Template{}, errors.New("test error")),
		)

		err := safednsTemplateCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new template: test error", err.Error())
	})
}

func Test_safednsTemplateUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := safednsTemplateUpdateCmd(nil).Args(nil, []string{"testdomain.com"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := safednsTemplateUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_safednsTemplateUpdate(t *testing.T) {
	t.Run("UpdateSingle_ByID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateUpdateCmd(nil)
		cmd.Flags().Set("default", "false")

		gomock.InOrder(
			service.EXPECT().PatchTemplate(123, gomock.Any()).Return(123, nil).Do(func(templateID int, patch safedns.PatchTemplateRequest) {
				if patch.Default == nil {
					t.Fatal("Expected non-nil default")
				}
				if *patch.Default != false {
					t.Fatalf("Expected default of false, got %t", *patch.Default)
				}
			}),
			service.EXPECT().GetTemplate(123).Return(safedns.Template{}, nil),
		)

		safednsTemplateUpdate(service, cmd, []string{"123"})
	})

	t.Run("UpdateMultiple_ByID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateUpdateCmd(nil)
		cmd.Flags().Set("default", "false")

		gomock.InOrder(
			service.EXPECT().PatchTemplate(123, gomock.Any()).Return(123, nil).Do(func(templateID int, patch safedns.PatchTemplateRequest) {
				if patch.Default == nil {
					t.Fatal("Expected non-nil default")
				}
				if *patch.Default != false {
					t.Fatalf("Expected default of false, got %t", *patch.Default)
				}
			}),
			service.EXPECT().GetTemplate(123).Return(safedns.Template{}, nil),
			service.EXPECT().PatchTemplate(456, gomock.Any()).Return(456, nil).Do(func(templateID int, patch safedns.PatchTemplateRequest) {
				if patch.Default == nil {
					t.Fatal("Expected non-nil default")
				}
				if *patch.Default != false {
					t.Fatalf("Expected default of false, got %t", *patch.Default)
				}
			}),
			service.EXPECT().GetTemplate(456).Return(safedns.Template{}, nil),
		)

		safednsTemplateUpdate(service, cmd, []string{"123", "456"})
	})

	t.Run("ErrorUpdatingTemplate_ByID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateUpdateCmd(nil)
		cmd.Flags().Set("default", "false")

		service.EXPECT().PatchTemplate(123, gomock.Any()).Return(0, errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error updating template [123]: test error\n", func() {
			safednsTemplateUpdate(service, cmd, []string{"123"})
		})
	})

	t.Run("UpdateSingle_ByName", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateUpdateCmd(nil)
		cmd.Flags().Set("default", "false")

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
			service.EXPECT().PatchTemplate(123, gomock.Any()).Return(123, nil).Do(func(templateID int, patch safedns.PatchTemplateRequest) {
				if patch.Default == nil {
					t.Fatal("Expected non-nil default")
				}
				if *patch.Default != false {
					t.Fatalf("Expected default of false, got %t", *patch.Default)
				}
			}),
			service.EXPECT().GetTemplate(123).Return(safedns.Template{}, nil),
		)

		safednsTemplateUpdate(service, cmd, []string{"test template 1"})
	})

	t.Run("UpdateMultiple_ByName", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateUpdateCmd(nil)
		cmd.Flags().Set("default", "false")

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
			service.EXPECT().PatchTemplate(123, gomock.Any()).Return(123, nil).Do(func(templateID int, patch safedns.PatchTemplateRequest) {
				if patch.Default == nil {
					t.Fatal("Expected non-nil default")
				}
				if *patch.Default != false {
					t.Fatalf("Expected default of false, got %t", *patch.Default)
				}
			}),
			service.EXPECT().GetTemplate(123).Return(safedns.Template{}, nil),

			service.EXPECT().GetTemplates(gomock.Eq(
				connection.APIRequestParameters{
					Filtering: []connection.APIRequestFiltering{
						connection.APIRequestFiltering{
							Property: "name",
							Operator: connection.EQOperator,
							Value:    []string{"test template 2"},
						},
					},
				}),
			).Return([]safedns.Template{safedns.Template{ID: 456}}, nil),
			service.EXPECT().PatchTemplate(456, gomock.Any()).Return(456, nil).Do(func(templateID int, patch safedns.PatchTemplateRequest) {
				if patch.Default == nil {
					t.Fatal("Expected non-nil default")
				}
				if *patch.Default != false {
					t.Fatalf("Expected default of false, got %t", *patch.Default)
				}
			}),
			service.EXPECT().GetTemplate(456).Return(safedns.Template{}, nil),
		)

		safednsTemplateUpdate(service, cmd, []string{"test template 1", "test template 2"})
	})

	t.Run("ErrorUpdatingTemplate_ByName_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateUpdateCmd(nil)
		cmd.Flags().Set("name", "new test template 1")

		service.EXPECT().PatchTemplate(123, gomock.Any()).Return(0, errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error updating template [123]: test error\n", func() {
			safednsTemplateUpdate(service, cmd, []string{"123"})
		})
	})

	t.Run("GetTemplatesError_OutputsError", func(t *testing.T) {
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
			).Return([]safedns.Template{safedns.Template{}}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error locating template [test template 1]: Error retrieving items: test error\n", func() {
			safednsTemplateUpdate(service, &cobra.Command{}, []string{"test template 1"})
		})
	})

	t.Run("GetTemplateError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateUpdateCmd(nil)
		cmd.Flags().Set("name", "new test template 1")

		gomock.InOrder(
			service.EXPECT().PatchTemplate(123, gomock.Any()).Return(123, nil),
			service.EXPECT().GetTemplate(123).Return(safedns.Template{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated template: test error\n", func() {
			safednsTemplateUpdate(service, cmd, []string{"123"})
		})
	})
}

func Test_safednsTemplateDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := safednsTemplateDeleteCmd(nil).Args(nil, []string{"test template 1"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := safednsTemplateDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_safednsTemplateDelete(t *testing.T) {
	t.Run("SingleTemplate_ByID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().DeleteTemplate(123).Return(nil).Times(1)

		safednsTemplateDelete(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleTemplates_ByID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteTemplate(123).Return(nil),
			service.EXPECT().DeleteTemplate(456).Return(nil),
		)

		safednsTemplateDelete(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("SingleTemplate_ByName", func(t *testing.T) {
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
			service.EXPECT().DeleteTemplate(123).Return(nil),
		)

		safednsTemplateDelete(service, &cobra.Command{}, []string{"test template 1"})
	})

	t.Run("MultipleTemplates_ByName", func(t *testing.T) {
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
			service.EXPECT().DeleteTemplate(123).Return(nil),
			service.EXPECT().GetTemplates(gomock.Eq(
				connection.APIRequestParameters{
					Filtering: []connection.APIRequestFiltering{
						connection.APIRequestFiltering{
							Property: "name",
							Operator: connection.EQOperator,
							Value:    []string{"test template 2"},
						},
					},
				}),
			).Return([]safedns.Template{safedns.Template{ID: 456}}, nil),
			service.EXPECT().DeleteTemplate(456).Return(nil),
		)

		safednsTemplateDelete(service, &cobra.Command{}, []string{"test template 1", "test template 2"})
	})

	t.Run("GetTemplatesError_OutputsError", func(t *testing.T) {
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
			).Return([]safedns.Template{safedns.Template{}}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error locating template [test template 1]: Error retrieving items: test error\n", func() {
			safednsTemplateDelete(service, &cobra.Command{}, []string{"test template 1"})
		})
	})

	t.Run("DeleteTemplateByNameError_OutputsError", func(t *testing.T) {
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
			service.EXPECT().DeleteTemplate(123).Return(errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error removing template [test template 1]: test error\n", func() {
			safednsTemplateDelete(service, &cobra.Command{}, []string{"test template 1"})
		})
	})

	t.Run("DeleteTemplateByIDError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().DeleteTemplate(123).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing template [123]: test error\n", func() {
			safednsTemplateDelete(service, &cobra.Command{}, []string{"123"})
		})
	})
}
