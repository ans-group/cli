package sharedexchange

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/sharedexchange"
)

func Test_sharedexchangeDomainList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSharedExchangeService(mockCtrl)

		service.EXPECT().GetDomains(gomock.Any()).Return([]sharedexchange.Domain{}, nil).Times(1)

		sharedexchangeDomainList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSharedExchangeService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := sharedexchangeDomainList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetDomainsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSharedExchangeService(mockCtrl)

		service.EXPECT().GetDomains(gomock.Any()).Return([]sharedexchange.Domain{}, errors.New("test error")).Times(1)

		err := sharedexchangeDomainList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving domains: test error", err.Error())
	})
}

func Test_sharedexchangeDomainShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := sharedexchangeDomainShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := sharedexchangeDomainShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_sharedexchangeDomainShow(t *testing.T) {
	t.Run("SingleDomain", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSharedExchangeService(mockCtrl)

		service.EXPECT().GetDomain(123).Return(sharedexchange.Domain{}, nil).Times(1)

		sharedexchangeDomainShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleDomains", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSharedExchangeService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetDomain(123).Return(sharedexchange.Domain{}, nil),
			service.EXPECT().GetDomain(456).Return(sharedexchange.Domain{}, nil),
		)

		sharedexchangeDomainShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetDomainID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSharedExchangeService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid domain ID [abc]\n", func() {
			sharedexchangeDomainShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetDomainError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSharedExchangeService(mockCtrl)

		service.EXPECT().GetDomain(123).Return(sharedexchange.Domain{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving domain [123]: test error\n", func() {
			sharedexchangeDomainShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
