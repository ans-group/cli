package loadbalancer

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_loadbalancerClusterACLTemplateShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerClusterACLTemplateShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerClusterACLTemplateShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing cluster", err.Error())
	})
}

func Test_loadbalancerClusterACLTemplateShow(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetClusterACLTemplates(123).Return(loadbalancer.ACLTemplates{}, nil).Times(1)

		loadbalancerClusterACLTemplateShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("GetClusterACLTemplatesError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetClusterACLTemplates(123).Return(loadbalancer.ACLTemplates{}, errors.New("test error")).Times(1)

		err := loadbalancerClusterACLTemplateShow(service, &cobra.Command{}, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "error retrieving ACL templates: test error", err.Error())
	})
}
