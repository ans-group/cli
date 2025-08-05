package ddosx

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ddosxSSLList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetSSLs(gomock.Any()).Return([]ddosx.SSL{}, nil).Times(1)

		ddosxSSLList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ddosxSSLList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSSLsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetSSLs(gomock.Any()).Return([]ddosx.SSL{}, errors.New("test error")).Times(1)

		err := ddosxSSLList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving ssls: test error", err.Error())
	})
}

func Test_ddosxSSLShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxSSLShowCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxSSLShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing ssl", err.Error())
	})
}

func Test_ddosxSSLShow(t *testing.T) {
	t.Run("SingleSSL", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetSSL("00000000-0000-0000-0000-000000000000").Return(ddosx.SSL{}, nil).Times(1)

		ddosxSSLShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MultipleSSLs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetSSL("00000000-0000-0000-0000-000000000000").Return(ddosx.SSL{}, nil),
			service.EXPECT().GetSSL("00000000-0000-0000-0000-000000000001").Return(ddosx.SSL{}, nil),
		)

		ddosxSSLShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("GetSSLError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetSSL("00000000-0000-0000-0000-000000000000").Return(ddosx.SSL{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving ssl [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxSSLShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxSSLCreate(t *testing.T) {
	t.Run("Valid_FromUKFastSSLID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLCreateCmd(nil, nil)
		cmd.Flags().Set("friendly-name", "testssl1")
		cmd.Flags().Set("ans-ssl-id", "123")

		expectedRequest := ddosx.CreateSSLRequest{
			FriendlyName: "testssl1",
			UKFastSSLID:  123,
		}

		gomock.InOrder(
			service.EXPECT().CreateSSL(gomock.Eq(expectedRequest)).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetSSL("00000000-0000-0000-0000-000000000000").Return(ddosx.SSL{}, nil),
		)

		ddosxSSLCreate(service, cmd, nil, []string{})
	})

	t.Run("Valid_FromCertificate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLCreateCmd(nil, nil)
		cmd.Flags().Set("friendly-name", "testssl1")
		cmd.Flags().Set("key", "testkey1")
		cmd.Flags().Set("certificate", "testcertificate1")
		cmd.Flags().Set("ca-bundle", "testcabundle1")

		expectedRequest := ddosx.CreateSSLRequest{
			FriendlyName: "testssl1",
			Key:          "testkey1",
			Certificate:  "testcertificate1",
			CABundle:     "testcabundle1",
		}

		gomock.InOrder(
			service.EXPECT().CreateSSL(gomock.Eq(expectedRequest)).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetSSL("00000000-0000-0000-0000-000000000000").Return(ddosx.SSL{}, nil),
		)

		ddosxSSLCreate(service, cmd, nil, []string{})
	})

	t.Run("CreateSSLError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLCreateCmd(nil, nil)
		cmd.Flags().Set("friendly-name", "testssl1")

		gomock.InOrder(
			service.EXPECT().CreateSSL(gomock.Any()).Return("00000000-0000-0000-0000-000000000000", errors.New("test error")),
		)

		err := ddosxSSLCreate(service, cmd, nil, []string{})
		assert.Equal(t, "error creating ssl: test error", err.Error())
	})

	t.Run("GetSSLError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLCreateCmd(nil, nil)
		cmd.Flags().Set("friendly-name", "testssl1")

		gomock.InOrder(
			service.EXPECT().CreateSSL(gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetSSL("00000000-0000-0000-0000-000000000000").Return(ddosx.SSL{}, errors.New("test error")),
		)

		err := ddosxSSLCreate(service, cmd, nil, []string{})
		assert.Equal(t, "error retrieving new ssl [00000000-0000-0000-0000-000000000000]: test error", err.Error())
	})
}

func Test_ddosxSSLUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxSSLUpdateCmd(nil, nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxSSLUpdateCmd(nil, nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing ssl", err.Error())
	})
}

func Test_ddosxSSLUpdate(t *testing.T) {
	t.Run("Valid_UKFastSSLID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLUpdateCmd(nil, nil)
		cmd.Flags().Set("friendly-name", "testssl1")
		cmd.Flags().Set("ans-ssl-id", "123")

		expectedRequest := ddosx.PatchSSLRequest{
			FriendlyName: "testssl1",
			UKFastSSLID:  123,
		}

		gomock.InOrder(
			service.EXPECT().PatchSSL("00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetSSL("00000000-0000-0000-0000-000000000000").Return(ddosx.SSL{}, nil),
		)

		ddosxSSLUpdate(service, cmd, nil, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("Valid_Certificate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLUpdateCmd(nil, nil)
		cmd.Flags().Set("friendly-name", "testssl1")
		cmd.Flags().Set("key", "testkey1")
		cmd.Flags().Set("certificate", "testcertificate1")
		cmd.Flags().Set("ca-bundle", "testcabundle1")

		expectedRequest := ddosx.PatchSSLRequest{
			FriendlyName: "testssl1",
			Key:          "testkey1",
			Certificate:  "testcertificate1",
			CABundle:     "testcabundle1",
		}

		gomock.InOrder(
			service.EXPECT().PatchSSL("00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetSSL("00000000-0000-0000-0000-000000000000").Return(ddosx.SSL{}, nil),
		)

		ddosxSSLUpdate(service, cmd, nil, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("UpdateSSLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLUpdateCmd(nil, nil)
		cmd.Flags().Set("friendly-name", "testssl1")

		gomock.InOrder(
			service.EXPECT().PatchSSL("00000000-0000-0000-0000-000000000000", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", errors.New("test error")),
		)

		err := ddosxSSLUpdate(service, cmd, nil, []string{"00000000-0000-0000-0000-000000000000"})
		assert.Equal(t, "error updating ssl: test error", err.Error())
	})

	t.Run("GetSSLError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLUpdateCmd(nil, nil)
		cmd.Flags().Set("friendly-name", "testssl1")

		gomock.InOrder(
			service.EXPECT().PatchSSL("00000000-0000-0000-0000-000000000000", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetSSL("00000000-0000-0000-0000-000000000000").Return(ddosx.SSL{}, errors.New("test error")),
		)

		err := ddosxSSLUpdate(service, cmd, nil, []string{"00000000-0000-0000-0000-000000000000"})
		assert.Equal(t, "error retrieving updated ssl: test error", err.Error())
	})
}

func Test_ddosxSSLDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxSSLDeleteCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxSSLDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing ssl", err.Error())
	})
}

func Test_ddosxSSLDelete(t *testing.T) {
	t.Run("SingleSSL", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteSSL("00000000-0000-0000-0000-000000000000").Return(nil).Times(1)

		ddosxSSLDelete(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MultipleSSLs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteSSL("00000000-0000-0000-0000-000000000000").Return(nil),
			service.EXPECT().DeleteSSL("00000000-0000-0000-0000-000000000001").Return(nil),
		)

		ddosxSSLDelete(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("DeleteSSLError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteSSL("00000000-0000-0000-0000-000000000000").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing ssl [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxSSLDelete(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}
