package cloudflare

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/cloudflare"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_cloudflareZoneList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)

		service.EXPECT().GetZones(gomock.Any()).Return([]cloudflare.Zone{}, nil).Times(1)

		cloudflareZoneList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := cloudflareZoneList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetZonesError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)

		service.EXPECT().GetZones(gomock.Any()).Return([]cloudflare.Zone{}, errors.New("test error")).Times(1)

		err := cloudflareZoneList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving zones: test error", err.Error())
	})
}

func Test_cloudflareZoneShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := cloudflareZoneShowCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := cloudflareZoneShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})
}

func Test_cloudflareZoneShow(t *testing.T) {
	t.Run("SingleZone", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)

		service.EXPECT().GetZone("00000000-0000-0000-0000-000000000000").Return(cloudflare.Zone{}, nil).Times(1)

		cloudflareZoneShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("GetZoneError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)

		service.EXPECT().GetZone("00000000-0000-0000-0000-000000000000").Return(cloudflare.Zone{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving zone [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			cloudflareZoneShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_cloudflareZoneCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)
		cmd := cloudflareZoneCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testzone"})

		req := cloudflare.CreateZoneRequest{
			Name: "testzone",
		}

		gomock.InOrder(
			service.EXPECT().CreateZone(req).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetZone("00000000-0000-0000-0000-000000000000").Return(cloudflare.Zone{}, nil),
		)

		cloudflareZoneCreate(service, cmd, []string{})
	})

	t.Run("CreateZoneError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)
		cmd := cloudflareZoneCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testzone"})

		service.EXPECT().CreateZone(gomock.Any()).Return("", errors.New("test error"))

		err := cloudflareZoneCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating zone: test error", err.Error())
	})

	t.Run("GetZoneError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)
		cmd := cloudflareZoneCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testzone"})

		gomock.InOrder(
			service.EXPECT().CreateZone(gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetZone("00000000-0000-0000-0000-000000000000").Return(cloudflare.Zone{}, errors.New("test error")),
		)

		err := cloudflareZoneCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new zone: test error", err.Error())
	})
}

func Test_cloudflareZoneUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := cloudflareZoneUpdateCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := cloudflareZoneUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})
}

func Test_cloudflareZoneUpdate(t *testing.T) {
	t.Run("SingleZone", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)

		cmd := cloudflareZoneUpdateCmd(nil)
		cmd.Flags().Set("subscription", "00000000-0000-0000-0000-000000000000")

		req := cloudflare.PatchZoneRequest{
			SubscriptionID: "00000000-0000-0000-0000-000000000000",
		}

		service.EXPECT().PatchZone("00000000-0000-0000-0000-000000000000", req).Return(nil).Times(1)

		cloudflareZoneUpdate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("UpdateZoneError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)

		service.EXPECT().PatchZone("00000000-0000-0000-0000-000000000000", cloudflare.PatchZoneRequest{}).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating zone [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			cloudflareZoneUpdate(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_cloudflareZoneDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := cloudflareZoneDeleteCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := cloudflareZoneDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})
}

func Test_cloudflareZoneDelete(t *testing.T) {
	t.Run("SingleZone", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)

		service.EXPECT().DeleteZone("00000000-0000-0000-0000-000000000000").Return(nil).Times(1)

		cloudflareZoneDelete(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("DeleteZoneError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)

		service.EXPECT().DeleteZone("00000000-0000-0000-0000-000000000000").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing zone [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			cloudflareZoneDelete(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}
