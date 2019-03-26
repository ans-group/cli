package cmd

import (
	"errors"
	"testing"

	"github.com/ukfast/sdk-go/pkg/connection"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func Test_ddosxDomainPropertyListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainPropertyListCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainPropertyListCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainPropertyList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainProperties("testdomain1.co.uk", gomock.Any()).Return([]ddosx.DomainProperty{}, nil).Times(1)

		ddosxDomainPropertyList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("WithFilter_AppendsFilter", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainPropertyListCmd()
		cmd.Flags().Set("name", "testproperty1")

		expectedParams := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				connection.APIRequestFiltering{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{"testproperty1"},
				},
			},
		}

		service.EXPECT().GetDomainProperties("testdomain1.co.uk", gomock.Eq(expectedParams)).Return([]ddosx.DomainProperty{}, nil).Times(1)

		ddosxDomainPropertyList(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		test_output.AssertFatalOutput(t, "Missing value for filtering\n", func() {
			ddosxDomainPropertyList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainPropertiesError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainProperties("testdomain1.co.uk", gomock.Any()).Return([]ddosx.DomainProperty{}, errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving domain properties: test error\n", func() {
			ddosxDomainPropertyList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainPropertyShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainPropertyShowCmd().Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		err := ddosxDomainPropertyShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingDomainProperty_Error", func(t *testing.T) {
		err := ddosxDomainPropertyShowCmd().Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain property", err.Error())
	})
}

func Test_ddosxDomainPropertyShow(t *testing.T) {
	t.Run("SingleDomain", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainProperty("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.DomainProperty{}, nil).Times(1)

		ddosxDomainPropertyShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MultipleDomains", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetDomainProperty("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.DomainProperty{}, nil),
			service.EXPECT().GetDomainProperty("testdomain1.co.uk", "00000000-0000-0000-0000-000000000001").Return(ddosx.DomainProperty{}, nil),
		)

		ddosxDomainPropertyShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("GetDomainPropertyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainProperty("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.DomainProperty{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving domain property [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainPropertyShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxDomainPropertyUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainPropertyUpdateCmd().Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		err := ddosxDomainPropertyUpdateCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingDomainProperty_Error", func(t *testing.T) {
		err := ddosxDomainPropertyUpdateCmd().Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain property", err.Error())
	})
}

func Test_ddosxDomainPropertyUpdate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainPropertyUpdateCmd()
		cmd.Flags().Set("value", "testvalue1")

		gomock.InOrder(
			service.EXPECT().PatchDomainProperty("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Do(func(domainName string, propertyUUID string, req ddosx.PatchDomainPropertyRequest) {
				if req.Value != "testvalue1" {
					t.Fatal("Expected value of testvalue1")
				}
			}).Return(nil),
			service.EXPECT().GetDomainProperty("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.DomainProperty{}, nil),
		)

		ddosxDomainPropertyUpdate(service, cmd, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("PatchDomainPropertyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().PatchDomainProperty("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating domain property [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainPropertyUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})

	t.Run("GetDomainPropertyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchDomainProperty("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(nil),
			service.EXPECT().GetDomainProperty("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.DomainProperty{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated domain property [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainPropertyUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}
