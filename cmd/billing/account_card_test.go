package billing

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/billing"
)

func Test_billingCardList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetCards(gomock.Any()).Return([]billing.Card{}, nil).Times(1)

		billingCardList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := billingCardList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetCardsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetCards(gomock.Any()).Return([]billing.Card{}, errors.New("test error")).Times(1)

		err := billingCardList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving cards: test error", err.Error())
	})
}

func Test_billingCardShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := billingCardShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := billingCardShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing card", err.Error())
	})
}

func Test_billingCardShow(t *testing.T) {
	t.Run("SingleCard", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetCard(123).Return(billing.Card{}, nil).Times(1)

		billingCardShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleCards", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetCard(123).Return(billing.Card{}, nil),
			service.EXPECT().GetCard(456).Return(billing.Card{}, nil),
		)

		billingCardShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetCardID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid card ID [abc]\n", func() {
			billingCardShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetCardError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().GetCard(123).Return(billing.Card{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving card [123]: test error\n", func() {
			billingCardShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_billingCardCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)
		cmd := billingCardCreateCmd(nil)
		cmd.Flags().Set("name", "test card 1")

		expectedReq := billing.CreateCardRequest{
			Name: "test card 1",
		}

		gomock.InOrder(
			service.EXPECT().CreateCard(expectedReq).Return(123, nil),
			service.EXPECT().GetCard(123).Return(billing.Card{}, nil),
		)

		billingCardCreate(service, cmd, []string{})
	})

	t.Run("CreateCardError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().CreateCard(gomock.Any()).Return(0, errors.New("test error")).Times(1)

		err := billingCardCreate(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error creating card: test error", err.Error())
	})

	t.Run("GetCardError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().CreateCard(gomock.Any()).Return(123, nil),
			service.EXPECT().GetCard(123).Return(billing.Card{}, errors.New("test error")),
		)

		err := billingCardCreate(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving new card: test error", err.Error())
	})
}

func Test_billingCardUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := billingCardUpdateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := billingCardUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing card", err.Error())
	})
}

func Test_billingCardUpdate(t *testing.T) {
	t.Run("SingleCard_SetsCPU", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)
		cmd := billingCardUpdateCmd(nil)
		cmd.Flags().Set("card-type", "visa")

		expectedPatch := billing.PatchCardRequest{
			CardType: "visa",
		}

		gomock.InOrder(
			service.EXPECT().PatchCard(123, expectedPatch).Return(nil),
			service.EXPECT().GetCard(123).Return(billing.Card{}, nil),
		)

		billingCardUpdate(service, cmd, []string{"123"})
	})

	t.Run("MultipleCards", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)
		cmd := billingCardUpdateCmd(nil)
		cmd.Flags().Set("card-type", "visa")

		expectedPatch := billing.PatchCardRequest{
			CardType: "visa",
		}

		gomock.InOrder(
			service.EXPECT().PatchCard(123, gomock.Eq(expectedPatch)).Return(nil),
			service.EXPECT().GetCard(123).Return(billing.Card{}, nil),
			service.EXPECT().PatchCard(456, gomock.Eq(expectedPatch)).Return(nil),
			service.EXPECT().GetCard(456).Return(billing.Card{}, nil),
		)

		billingCardUpdate(service, cmd, []string{"123", "456"})
	})

	t.Run("InvalidCardID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid card ID [abc]\n", func() {
			billingCardUpdate(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("PatchCardError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().PatchCard(123, gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating card [123]: test error\n", func() {
			billingCardUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})

	t.Run("GetCardError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchCard(123, gomock.Any()).Return(nil),
			service.EXPECT().GetCard(123).Return(billing.Card{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated card [123]: test error\n", func() {
			billingCardUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_billingCardDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := billingCardDeleteCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := billingCardDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing card", err.Error())
	})
}

func Test_billingCardDelete(t *testing.T) {
	t.Run("SingleCard", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().DeleteCard(123).Return(nil).Times(1)

		billingCardDelete(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleCards", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteCard(123).Return(nil),
			service.EXPECT().DeleteCard(456).Return(nil),
		)

		billingCardDelete(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("InvalidCardID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid card ID [abc]\n", func() {
			billingCardDelete(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetCardError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockBillingService(mockCtrl)

		service.EXPECT().DeleteCard(123).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing card [123]: test error\n", func() {
			billingCardDelete(service, &cobra.Command{}, []string{"123"})
		})
	})
}
