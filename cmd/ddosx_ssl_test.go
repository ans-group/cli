package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func Test_ddosxSSLList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetSSLs(gomock.Any()).Return([]ddosx.SSL{}, nil).Times(1)

		ddosxSSLList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		test_output.AssertFatalOutput(t, "Missing value for filtering\n", func() {
			ddosxSSLList(service, cmd, []string{})
		})
	})

	t.Run("GetSSLsError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetSSLs(gomock.Any()).Return([]ddosx.SSL{}, errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving ssls: test error\n", func() {
			ddosxSSLList(service, &cobra.Command{}, []string{})
		})
	})
}

func Test_ddosxSSLShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxSSLShowCmd().Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxSSLShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ssl", err.Error())
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
		cmd := ddosxSSLCreateCmd()
		cmd.Flags().Set("friendly-name", "testssl1")
		cmd.Flags().Set("ukfast-ssl-id", "123")

		expectedRequest := ddosx.CreateSSLRequest{
			FriendlyName: "testssl1",
			UKFastSSLID:  123,
		}

		gomock.InOrder(
			service.EXPECT().CreateSSL(gomock.Eq(expectedRequest)).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetSSL("00000000-0000-0000-0000-000000000000").Return(ddosx.SSL{}, nil),
		)

		ddosxSSLCreate(service, cmd, []string{})
	})

	t.Run("Valid_FromCertificate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLCreateCmd()
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

		ddosxSSLCreate(service, cmd, []string{})
	})

	t.Run("CreateSSLError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLCreateCmd()
		cmd.Flags().Set("friendly-name", "testssl1")

		gomock.InOrder(
			service.EXPECT().CreateSSL(gomock.Any()).Return("00000000-0000-0000-0000-000000000000", errors.New("test error")),
		)

		test_output.AssertFatalOutput(t, "Error creating ssl: test error\n", func() {
			ddosxSSLCreate(service, cmd, []string{})
		})
	})

	t.Run("GetSSLError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLCreateCmd()
		cmd.Flags().Set("friendly-name", "testssl1")

		gomock.InOrder(
			service.EXPECT().CreateSSL(gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetSSL("00000000-0000-0000-0000-000000000000").Return(ddosx.SSL{}, errors.New("test error")),
		)

		test_output.AssertFatalOutput(t, "Error retrieving new ssl [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxSSLCreate(service, cmd, []string{})
		})
	})
}

func Test_ddosxSSLUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxSSLUpdateCmd().Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxSSLUpdateCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ssl", err.Error())
	})
}

func Test_ddosxSSLUpdate(t *testing.T) {
	t.Run("Valid_UKFastSSLID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLUpdateCmd()
		cmd.Flags().Set("friendly-name", "testssl1")
		cmd.Flags().Set("ukfast-ssl-id", "123")

		expectedRequest := ddosx.PatchSSLRequest{
			FriendlyName: "testssl1",
			UKFastSSLID:  123,
		}

		gomock.InOrder(
			service.EXPECT().PatchSSL("00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetSSL("00000000-0000-0000-0000-000000000000").Return(ddosx.SSL{}, nil),
		)

		ddosxSSLUpdate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("Valid_Certificate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLUpdateCmd()
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

		ddosxSSLUpdate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("UpdateSSLError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLUpdateCmd()
		cmd.Flags().Set("friendly-name", "testssl1")

		gomock.InOrder(
			service.EXPECT().PatchSSL("00000000-0000-0000-0000-000000000000", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", errors.New("test error")),
		)

		test_output.AssertFatalOutput(t, "Error updating ssl: test error\n", func() {
			ddosxSSLUpdate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})

	t.Run("GetSSLError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxSSLUpdateCmd()
		cmd.Flags().Set("friendly-name", "testssl1")

		gomock.InOrder(
			service.EXPECT().PatchSSL("00000000-0000-0000-0000-000000000000", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetSSL("00000000-0000-0000-0000-000000000000").Return(ddosx.SSL{}, errors.New("test error")),
		)

		test_output.AssertFatalOutput(t, "Error retrieving updated ssl: test error\n", func() {
			ddosxSSLUpdate(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxSSLDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxSSLDeleteCmd().Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxSSLDeleteCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing ssl", err.Error())
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
