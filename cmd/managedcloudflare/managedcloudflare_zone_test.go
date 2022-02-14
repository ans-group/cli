package managedcloudflare

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

func Test_managedcloudflareZoneList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().GetZones(gomock.Any()).Return([]managedcloudflare.Zone{}, nil).Times(1)

		managedcloudflareZoneList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := managedcloudflareZoneList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetZonesError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().GetZones(gomock.Any()).Return([]managedcloudflare.Zone{}, errors.New("test error")).Times(1)

		err := managedcloudflareZoneList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving zones: test error", err.Error())
	})
}

func Test_managedcloudflareZoneShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := managedcloudflareZoneShowCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := managedcloudflareZoneShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})
}

func Test_managedcloudflareZoneShow(t *testing.T) {
	t.Run("SingleZone", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().GetZone("00000000-0000-0000-0000-000000000000").Return(managedcloudflare.Zone{}, nil).Times(1)

		managedcloudflareZoneShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("GetZoneError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().GetZone("00000000-0000-0000-0000-000000000000").Return(managedcloudflare.Zone{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving zone [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			managedcloudflareZoneShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_managedcloudflareZoneCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)
		cmd := managedcloudflareZoneCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testzone"})

		req := managedcloudflare.CreateZoneRequest{
			Name: "testzone",
		}

		gomock.InOrder(
			service.EXPECT().CreateZone(req).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetZone("00000000-0000-0000-0000-000000000000").Return(managedcloudflare.Zone{}, nil),
		)

		managedcloudflareZoneCreate(service, cmd, []string{})
	})

	t.Run("CreateZoneError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)
		cmd := managedcloudflareZoneCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testzone"})

		service.EXPECT().CreateZone(gomock.Any()).Return("", errors.New("test error"))

		err := managedcloudflareZoneCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating zone: test error", err.Error())
	})

	t.Run("GetZoneError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)
		cmd := managedcloudflareZoneCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testzone"})

		gomock.InOrder(
			service.EXPECT().CreateZone(gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetZone("00000000-0000-0000-0000-000000000000").Return(managedcloudflare.Zone{}, errors.New("test error")),
		)

		err := managedcloudflareZoneCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new zone: test error", err.Error())
	})
}

func Test_managedcloudflareZoneDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := managedcloudflareZoneDeleteCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := managedcloudflareZoneDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})
}

func Test_managedcloudflareZoneDelete(t *testing.T) {
	t.Run("SingleZone", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().DeleteZone("00000000-0000-0000-0000-000000000000").Return(nil).Times(1)

		managedcloudflareZoneDelete(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("DeleteZoneError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)

		service.EXPECT().DeleteZone("00000000-0000-0000-0000-000000000000").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing zone [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			managedcloudflareZoneDelete(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}
