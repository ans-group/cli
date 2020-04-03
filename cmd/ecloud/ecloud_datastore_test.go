package ecloud

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
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

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudDatastoreList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetDatastoresError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetDatastores(gomock.Any()).Return([]ecloud.Datastore{}, errors.New("test error")).Times(1)

		err := ecloudDatastoreList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving datastores: test error", err.Error())
	})
}

func Test_ecloudDatastoreShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudDatastoreShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudDatastoreShowCmd(nil).Args(nil, []string{})

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

		test_output.AssertErrorOutput(t, "Invalid datastore ID [abc]\n", func() {
			ecloudDatastoreShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetDatastoreError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetDatastore(123).Return(ecloud.Datastore{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving datastore [123]: test error\n", func() {
			ecloudDatastoreShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
