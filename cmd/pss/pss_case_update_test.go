package pss

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_pssCaseUpdateListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssCaseUpdateListCmd(nil).Args(nil, []string{"INC123456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := pssCaseUpdateListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing case", err.Error())
	})
}

func Test_pssCaseUpdateList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetCaseUpdates("CHG123456", gomock.Any()).Return([]pss.CaseUpdate{}, nil).Times(1)

		pssCaseUpdateList(service, &cobra.Command{}, []string{"CHG123456"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := pssCaseUpdateList(service, cmd, []string{"CHG123456"})
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetRequestsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetCaseUpdates("CHG123456", gomock.Any()).Return([]pss.CaseUpdate{}, errors.New("test error")).Times(1)

		err := pssCaseUpdateList(service, &cobra.Command{}, []string{"CHG123456"})
		assert.Equal(t, "error retrieving case updates: test error", err.Error())
	})
}

func Test_pssCaseUpdateShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := pssCaseUpdateShowCmd(nil).Args(nil, []string{"INC123456", "a9bf156e-bd33-4d85-b5f1-ed6eda21a72e"})

		assert.Nil(t, err)
	})

	t.Run("NoCaseID_Error", func(t *testing.T) {
		err := pssCaseUpdateShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing case", err.Error())
	})

	t.Run("NoCaseUpdateID_Error", func(t *testing.T) {
		err := pssCaseUpdateShowCmd(nil).Args(nil, []string{"INC123456"})

		assert.NotNil(t, err)
		assert.Equal(t, "missing case update", err.Error())
	})
}

func Test_pssCaseUpdateShow(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetCaseUpdate("CHG123456", "79b7aa31-3755-460f-9275-453053630959").Return(pss.CaseUpdate{}, nil).Times(1)

		pssCaseUpdateShow(service, &cobra.Command{}, []string{"CHG123456", "79b7aa31-3755-460f-9275-453053630959"})
	})

	t.Run("GetCaseUpdatesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockPSSService(mockCtrl)

		service.EXPECT().GetCaseUpdate("CHG123456", "a4928803-b5cd-4b55-9748-a54112d91769").Return(pss.CaseUpdate{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving case update [a4928803-b5cd-4b55-9748-a54112d91769]: test error\n", func() {
			pssCaseUpdateShow(service, &cobra.Command{}, []string{"CHG123456", "a4928803-b5cd-4b55-9748-a54112d91769"})
		})
	})
}
