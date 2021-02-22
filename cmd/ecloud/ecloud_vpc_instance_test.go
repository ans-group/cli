package ecloud

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudVPCInstanceListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPCInstanceListCmd(nil).Args(nil, []string{"vpc-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPCInstanceListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VPC", err.Error())
	})
}

func Test_ecloudVPCInstanceList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPCInstances("vpc-abcdef12", gomock.Any()).Return([]ecloud.Instance{}, nil).Times(1)

		ecloudVPCInstanceList(service, &cobra.Command{}, []string{"vpc-abcdef12"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVPCInstanceList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVPCsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPCInstances("vpc-abcdef12", gomock.Any()).Return([]ecloud.Instance{}, errors.New("test error")).Times(1)

		err := ecloudVPCInstanceList(service, &cobra.Command{}, []string{"vpc-abcdef12"})

		assert.Equal(t, "Error retrieving VPC instances: test error", err.Error())
	})
}
