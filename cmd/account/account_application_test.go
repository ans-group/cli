package account

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/account"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestValidateIPRanges(t *testing.T) {
	t.Run("ValidIPAddresses_ReturnsNil", func(t *testing.T) {
		ipRanges := []string{"192.168.1.1", "10.0.0.1", "172.16.0.1"}
		err := validateIPRanges(ipRanges)
		assert.Nil(t, err)
	})

	t.Run("ValidCIDRRanges_ReturnsNil", func(t *testing.T) {
		ipRanges := []string{"192.168.1.0/24", "10.0.0.0/8", "172.16.0.0/16"}
		err := validateIPRanges(ipRanges)
		assert.Nil(t, err)
	})

	t.Run("MixedValidIPsAndCIDR_ReturnsNil", func(t *testing.T) {
		ipRanges := []string{"192.168.1.1", "10.0.0.0/8", "172.16.0.1"}
		err := validateIPRanges(ipRanges)
		assert.Nil(t, err)
	})

	t.Run("InvalidIPAddress_ReturnsError", func(t *testing.T) {
		ipRanges := []string{"192.168.1.1", "invalid-ip", "172.16.0.1"}
		err := validateIPRanges(ipRanges)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid IP address or CIDR range: invalid-ip")
	})

	t.Run("InvalidCIDRRange_ReturnsError", func(t *testing.T) {
		ipRanges := []string{"192.168.1.0/24", "10.0.0.0/99", "172.16.0.0/16"}
		err := validateIPRanges(ipRanges)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid IP address or CIDR range: 10.0.0.0/99")
	})

	t.Run("EmptySlice_ReturnsNil", func(t *testing.T) {
		ipRanges := []string{}
		err := validateIPRanges(ipRanges)
		assert.Nil(t, err)
	})
}

func Test_accountApplicationList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetApplications(gomock.Any()).Return([]account.Application{}, nil).Times(1)

		accountApplicationList(service, &cobra.Command{}, []string{})
	})

	t.Run("GetApplicationsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetApplications(gomock.Any()).Return([]account.Application{}, errors.New("test error")).Times(1)

		err := accountApplicationList(service, &cobra.Command{}, []string{})

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "test error")
	})
}

func Test_accountApplicationShow(t *testing.T) {
	t.Run("SingleApplication", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetApplication("test-id").Return(account.Application{}, nil).Times(1)

		accountApplicationShow(service, &cobra.Command{}, []string{"test-id"})
	})

	t.Run("MultipleApplications", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().GetApplication("test-id1").Return(account.Application{}, nil).Times(1)
		service.EXPECT().GetApplication("test-id2").Return(account.Application{}, nil).Times(1)

		accountApplicationShow(service, &cobra.Command{}, []string{"test-id1", "test-id2"})
	})
}

func Test_accountApplicationDelete(t *testing.T) {
	t.Run("SingleApplication", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().DeleteApplication("test-id").Return(nil).Times(1)

		accountApplicationDelete(service, &cobra.Command{}, []string{"test-id"})
	})

	t.Run("MultipleApplications", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockAccountService(mockCtrl)

		service.EXPECT().DeleteApplication("test-id1").Return(nil).Times(1)
		service.EXPECT().DeleteApplication("test-id2").Return(nil).Times(1)

		accountApplicationDelete(service, &cobra.Command{}, []string{"test-id1", "test-id2"})
	})
}