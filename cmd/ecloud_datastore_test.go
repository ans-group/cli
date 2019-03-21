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
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudDatastoreList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetDatastores(gomock.Any()).Return([]ecloud.Datastore{}, nil).Times(1)

		ecloudDatastoreList(service, &cobra.Command{}, []string{})
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

		service := mocks.NewMockECloudService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		output := test.CatchStdErr(t, func() {
			ecloudDatastoreList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Missing value for filtering\n", output)
	})

	t.Run("GetDatastoresError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetDatastores(gomock.Any()).Return([]ecloud.Datastore{}, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			ecloudDatastoreList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving datastores: test error\n", output)
	})
}

func Test_ecloudDatastoreShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudDatastoreShowCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudDatastoreShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing datastore", err.Error())
	})
}

func Test_ecloudDatastoreShow(t *testing.T) {
	t.Run("SingleDatastore", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetDatastore(123).Return(ecloud.Datastore{}, nil).Times(1)

		ecloudDatastoreShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleDatastores", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetDatastore(123).Return(ecloud.Datastore{}, nil),
			service.EXPECT().GetDatastore(456).Return(ecloud.Datastore{}, nil),
		)

		ecloudDatastoreShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetDatastoreID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		output := test.CatchStdErr(t, func() {
			ecloudDatastoreShow(service, &cobra.Command{}, []string{"abc"})
		})

		assert.Equal(t, "Invalid datastore ID [abc]\n", output)
	})

	t.Run("GetDatastoreError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetDatastore(123).Return(ecloud.Datastore{}, errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			ecloudDatastoreShow(service, &cobra.Command{}, []string{"123"})
		})

		assert.Equal(t, "Error retrieving datastore [123]: test error\n", output)
	})
}
