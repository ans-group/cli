package ssl

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ssl"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_sslRecommendationsShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := sslRecommendationsShowCmd(nil).Args(nil, []string{"example.com"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := sslRecommendationsShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_sslRecommendationsShow(t *testing.T) {
	t.Run("SingleRecommendations", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetRecommendations("example.com").Return(ssl.Recommendations{}, nil).Times(1)

		sslRecommendationsShow(service, &cobra.Command{}, []string{"example.com"})
	})

	t.Run("MultipleRecommendations", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetRecommendations("example.com").Return(ssl.Recommendations{}, nil),
			service.EXPECT().GetRecommendations("example2.com").Return(ssl.Recommendations{}, nil),
		)

		sslRecommendationsShow(service, &cobra.Command{}, []string{"example.com", "example2.com"})
	})

	t.Run("GetRecommendationsError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetRecommendations("example.com").Return(ssl.Recommendations{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving SSL recommendations for domain [example.com]: test error\n", func() {
			sslRecommendationsShow(service, &cobra.Command{}, []string{"example.com"})
		})
	})
}
