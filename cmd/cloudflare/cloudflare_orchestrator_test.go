package cloudflare

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/cloudflare"
)

func Test_cloudflareOrchestratorCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)
		cmd := cloudflareOrchestratorCreateCmd(nil)
		cmd.ParseFlags([]string{"--zone-name=test"})

		req := cloudflare.CreateOrchestrationRequest{
			ZoneName: "test",
		}

		service.EXPECT().CreateOrchestration(req).Return(nil)

		cloudflareOrchestratorCreate(service, cmd, []string{})
	})

	t.Run("CreateOrchestratorError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockCloudflareService(mockCtrl)
		cmd := cloudflareOrchestratorCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testorchestrator"})

		service.EXPECT().CreateOrchestration(gomock.Any()).Return(errors.New("test error"))

		err := cloudflareOrchestratorCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating orchestration: test error", err.Error())
	})
}
