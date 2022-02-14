package managedcloudflare

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

func Test_managedcloudflareOrchestratorCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)
		cmd := managedcloudflareOrchestratorCreateCmd(nil)
		cmd.ParseFlags([]string{"--zone-name=test"})

		req := managedcloudflare.CreateOrchestrationRequest{
			ZoneName: "test",
		}

		service.EXPECT().CreateOrchestration(req).Return(nil)

		managedcloudflareOrchestratorCreate(service, cmd, []string{})
	})

	t.Run("CreateOrchestratorError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockManagedCloudflareService(mockCtrl)
		cmd := managedcloudflareOrchestratorCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testorchestrator"})

		service.EXPECT().CreateOrchestration(gomock.Any()).Return(errors.New("test error"))

		err := managedcloudflareOrchestratorCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating orchestration: test error", err.Error())
	})
}
