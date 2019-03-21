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
		cmd := safednsTemplateListCmd()
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
			safednsTemplateList(service, &cobra.Command{}, []string{"123"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Missing value for filtering\n", output)
	})

	t.Run("GetTemplatesError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetTemplates(gomock.Any()).Return([]safedns.Template{}, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			safednsTemplateList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving templates: test error\n", output)
	})
}

func Test_safednsTemplateShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := safednsTemplateShowCmd().Args(nil, []string{"test template 1"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := safednsTemplateShowCmd().Args(nil, []string{})

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

		output := test.CatchStdErr(t, func() {
			safednsTemplateShow(service, &cobra.Command{}, []string{"123"})
		})

		assert.Equal(t, "Error retrieving template [123]: Error retrieving template by ID [123]: test error\n", output)
	})
}

func Test_safednsTemplateCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateCreateCmd()
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

	t.Run("CreateTemplateError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})

		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateCreateCmd()
		cmd.Flags().Set("name", "test template 1")
		cmd.Flags().Set("default", "true")

		expectedRequest := safedns.CreateTemplateRequest{
			Name:    "test template 1",
			Default: true,
		}

		service.EXPECT().CreateTemplate(expectedRequest).Return(0, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			safednsTemplateCreate(service, cmd, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error creating template: test error\n", output)
	})

	t.Run("GetTemplateError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})

		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateCreateCmd()
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

		output := test.CatchStdErr(t, func() {
			safednsTemplateCreate(service, cmd, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving new template: test error\n", output)
	})
}

func Test_safednsTemplateUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := safednsTemplateUpdateCmd().Args(nil, []string{"testdomain.com"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := safednsTemplateUpdateCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing template", err.Error())
	})
}

func Test_safednsTemplateUpdate(t *testing.T) {
	t.Run("UpdateSingle_ByID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateUpdateCmd()
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
		cmd := safednsTemplateUpdateCmd()
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
		cmd := safednsTemplateUpdateCmd()
		cmd.Flags().Set("default", "false")

		service.EXPECT().PatchTemplate(123, gomock.Any()).Return(0, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			safednsTemplateUpdate(service, cmd, []string{"123"})
		})

		assert.Equal(t, "Error updating template [123]: test error\n", output)
	})

	t.Run("UpdateSingle_ByName", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateUpdateCmd()
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
		cmd := safednsTemplateUpdateCmd()
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
		cmd := safednsTemplateUpdateCmd()
		cmd.Flags().Set("name", "new test template 1")

		service.EXPECT().PatchTemplate(123, gomock.Any()).Return(0, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			safednsTemplateUpdate(service, cmd, []string{"123"})
		})

		assert.Equal(t, "Error updating template [123]: test error\n", output)
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

		output := test.CatchStdErr(t, func() {
			safednsTemplateUpdate(service, &cobra.Command{}, []string{"test template 1"})
		})

		assert.Equal(t, "Error locating template [test template 1]: Error retrieving items: test error\n", output)
	})

	t.Run("GetTemplateError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsTemplateUpdateCmd()
		cmd.Flags().Set("name", "new test template 1")

		gomock.InOrder(
			service.EXPECT().PatchTemplate(123, gomock.Any()).Return(123, nil),
			service.EXPECT().GetTemplate(123).Return(safedns.Template{}, errors.New("test error")),
		)

		output := test.CatchStdErr(t, func() {
			safednsTemplateUpdate(service, cmd, []string{"123"})
		})

		assert.Equal(t, "Error retrieving updated template: test error\n", output)
	})
}

func Test_safednsTemplateDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := safednsTemplateDeleteCmd().Args(nil, []string{"test template 1"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := safednsTemplateDeleteCmd().Args(nil, []string{})

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

		output := test.CatchStdErr(t, func() {
			safednsTemplateDelete(service, &cobra.Command{}, []string{"test template 1"})
		})

		assert.Equal(t, "Error locating template [test template 1]: Error retrieving items: test error\n", output)
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

		output := test.CatchStdErr(t, func() {
			safednsTemplateDelete(service, &cobra.Command{}, []string{"test template 1"})
		})

		assert.Equal(t, "Error removing template [test template 1]: test error\n", output)
	})

	t.Run("DeleteTemplateByIDError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().DeleteTemplate(123).Return(errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			safednsTemplateDelete(service, &cobra.Command{}, []string{"123"})
		})

		assert.Equal(t, "Error removing template [123]: test error\n", output)
	})
}
