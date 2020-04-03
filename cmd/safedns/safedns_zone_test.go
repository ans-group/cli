package safedns

import (
	"errors"
	"testing"

	"github.com/ukfast/sdk-go/pkg/connection"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
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
		cmd := safednsZoneListCmd(nil)
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

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := safednsZoneList(service, cmd, []string{"testdomain1.co.uk"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetZonesError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetZones(gomock.Any()).Return([]safedns.Zone{}, errors.New("test error")).Times(1)

		err := safednsZoneList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving zones: test error", err.Error())
	})
}

func Test_safednsZoneShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := safednsZoneShowCmd(nil).Args(nil, []string{"testdomain.com"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := safednsZoneShowCmd(nil).Args(nil, []string{})

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
		cmd := safednsZoneCreateCmd(nil)
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

	t.Run("CreateZoneError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsZoneCreateCmd(nil)
		cmd.Flags().Set("name", "testdomain1.com")
		cmd.Flags().Set("description", "test description")

		expectedRequest := safedns.CreateZoneRequest{
			Name:        "testdomain1.com",
			Description: "test description",
		}

		service.EXPECT().CreateZone(expectedRequest).Return(errors.New("test error")).Times(1)

		err := safednsZoneCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating zone: test error", err.Error())
	})

	t.Run("GetZoneError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsZoneCreateCmd(nil)
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

		err := safednsZoneCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new zone: test error", err.Error())
	})
}

func Test_safednsZoneUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := safednsZoneUpdateCmd(nil).Args(nil, []string{"testdomain.com"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := safednsZoneUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})
}

func Test_safednsZoneUpdate(t *testing.T) {
	t.Run("SingleZone", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		cmd := safednsZoneCreateCmd(nil)
		cmd.Flags().Set("description", "test description")

		expectedRequest := safedns.PatchZoneRequest{
			Description: "test description",
		}

		gomock.InOrder(
			service.EXPECT().PatchZone("testdomain1.com", expectedRequest).Return(nil),
			service.EXPECT().GetZone("testdomain1.com").Return(safedns.Zone{}, nil),
		)

		safednsZoneUpdate(service, cmd, []string{"testdomain1.com"})
	})

	t.Run("MultipleZones", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchZone("testdomain1.com", gomock.Any()).Return(nil),
			service.EXPECT().GetZone("testdomain1.com").Return(safedns.Zone{}, nil),
			service.EXPECT().PatchZone("testdomain2.com", gomock.Any()).Return(nil),
			service.EXPECT().GetZone("testdomain2.com").Return(safedns.Zone{}, nil),
		)

		safednsZoneUpdate(service, &cobra.Command{}, []string{"testdomain1.com", "testdomain2.com"})
	})

	t.Run("PatchZoneError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().PatchZone("testdomain1.com", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating zone [testdomain1.com]: test error\n", func() {
			safednsZoneUpdate(service, &cobra.Command{}, []string{"testdomain1.com"})
		})
	})

	t.Run("GetZoneError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchZone("testdomain1.com", gomock.Any()).Return(nil),
			service.EXPECT().GetZone("testdomain1.com").Return(safedns.Zone{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated zone [testdomain1.com]: test error\n", func() {
			safednsZoneUpdate(service, &cobra.Command{}, []string{"testdomain1.com"})
		})
	})
}

func Test_safednsZoneDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := safednsZoneDeleteCmd(nil).Args(nil, []string{"testdomain.com"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := safednsZoneDeleteCmd(nil).Args(nil, []string{})

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
