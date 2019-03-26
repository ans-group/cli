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
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func Test_safednsZoneList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetZones(gomock.Any()).Return([]safedns.Zone{}, nil).Times(1)

		safednsZoneList(service, &cobra.Command{}, []string{})
	})

	t.Run("ExpectedFilterFromFlags", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsZoneListCmd()
		cmd.Flags().Set("name", "testdomain1.co.uk")

		expectedParams := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				connection.APIRequestFiltering{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{"testdomain1.co.uk"},
				},
			},
		}

		service.EXPECT().GetZones(gomock.Eq(expectedParams)).Return([]safedns.Zone{}, nil).Times(1)

		safednsZoneList(service, cmd, []string{})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		test_output.AssertFatalOutput(t, "Missing value for filtering\n", func() {
			safednsZoneList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetZonesError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetZones(gomock.Any()).Return([]safedns.Zone{}, errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving zones: test error\n", func() {
			safednsZoneList(service, &cobra.Command{}, []string{})
		})
	})
}

func Test_safednsZoneShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := safednsZoneShowCmd().Args(nil, []string{"testdomain.com"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := safednsZoneShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})
}

func Test_safednsZoneShow(t *testing.T) {
	t.Run("SingleZone", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetZone("testdomain1.com").Return(safedns.Zone{}, nil).Times(1)

		safednsZoneShow(service, &cobra.Command{}, []string{"testdomain1.com"})
	})

	t.Run("MultipleZones", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetZone("testdomain1.com").Return(safedns.Zone{}, nil),
			service.EXPECT().GetZone("testdomain2.com").Return(safedns.Zone{}, nil),
		)

		safednsZoneShow(service, &cobra.Command{}, []string{"testdomain1.com", "testdomain2.com"})
	})

	t.Run("GetZoneError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetZone("testdomain1.com").Return(safedns.Zone{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving zone [testdomain1.com]: test error\n", func() {
			safednsZoneShow(service, &cobra.Command{}, []string{"testdomain1.com"})
		})
	})
}

func Test_safednsZoneCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsZoneCreateCmd()
		cmd.Flags().Set("name", "testdomain1.com")
		cmd.Flags().Set("description", "test description")

		expectedRequest := safedns.CreateZoneRequest{
			Name:        "testdomain1.com",
			Description: "test description",
		}

		gomock.InOrder(
			service.EXPECT().CreateZone(expectedRequest).Return(nil),
			service.EXPECT().GetZone("testdomain1.com").Return(safedns.Zone{}, nil),
		)

		safednsZoneCreate(service, cmd, []string{})
	})

	t.Run("CreateZoneError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsZoneCreateCmd()
		cmd.Flags().Set("name", "testdomain1.com")
		cmd.Flags().Set("description", "test description")

		expectedRequest := safedns.CreateZoneRequest{
			Name:        "testdomain1.com",
			Description: "test description",
		}

		service.EXPECT().CreateZone(expectedRequest).Return(errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error creating zone: test error\n", func() {
			safednsZoneCreate(service, cmd, []string{})
		})
	})

	t.Run("GetZoneError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsZoneCreateCmd()
		cmd.Flags().Set("name", "testdomain1.com")
		cmd.Flags().Set("description", "test description")

		expectedRequest := safedns.CreateZoneRequest{
			Name:        "testdomain1.com",
			Description: "test description",
		}

		gomock.InOrder(
			service.EXPECT().CreateZone(expectedRequest).Return(nil),
			service.EXPECT().GetZone("testdomain1.com").Return(safedns.Zone{}, errors.New("test error")),
		)

		test_output.AssertFatalOutput(t, "Error retrieving new zone: test error\n", func() {
			safednsZoneCreate(service, cmd, []string{})
		})
	})
}

func Test_safednsZoneDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := safednsZoneDeleteCmd().Args(nil, []string{"testdomain.com"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := safednsZoneDeleteCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})
}

func Test_safednsZoneDelete(t *testing.T) {
	t.Run("SingleZone", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().DeleteZone("testdomain1.com").Return(nil).Times(1)

		safednsZoneDelete(service, &cobra.Command{}, []string{"testdomain1.com"})
	})

	t.Run("MultipleZones", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteZone("testdomain1.com").Return(nil),
			service.EXPECT().DeleteZone("testdomain2.com").Return(nil),
		)

		safednsZoneDelete(service, &cobra.Command{}, []string{"testdomain1.com", "testdomain2.com"})
	})

	t.Run("DeleteZoneError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().DeleteZone("testdomain1.com").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing zone [testdomain1.com]: test error\n", func() {
			safednsZoneDelete(service, &cobra.Command{}, []string{"testdomain1.com"})
		})
	})
}
