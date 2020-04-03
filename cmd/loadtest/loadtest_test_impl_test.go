package loadtest

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func Test_loadtestTestList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetTests(gomock.Any()).Return([]ltaas.Test{}, nil).Times(1)

		loadtestTestList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadtestTestList(service, cmd, []string{})

		assert.NotNil(t, err)
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetTestsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetTests(gomock.Any()).Return([]ltaas.Test{}, errors.New("test error")).Times(1)

		err := loadtestTestList(service, &cobra.Command{}, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving tests: test error", err.Error())
	})
}

func Test_loadtestTestShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadtestTestShowCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadtestTestShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing test", err.Error())
	})
}

func Test_loadtestTestShow(t *testing.T) {
	t.Run("SingleTest", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetTest("00000000-0000-0000-0000-000000000000").Return(ltaas.Test{}, nil).Times(1)

		loadtestTestShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MultipleTests", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetTest("00000000-0000-0000-0000-000000000000").Return(ltaas.Test{}, nil),
			service.EXPECT().GetTest("00000000-0000-0000-0000-000000000001").Return(ltaas.Test{}, nil),
		)

		loadtestTestShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("GetTestError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().GetTest("00000000-0000-0000-0000-000000000000").Return(ltaas.Test{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving test [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			loadtestTestShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_loadtestTestDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadtestTestDeleteCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadtestTestDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing test", err.Error())
	})
}

func Test_loadtestTestDelete(t *testing.T) {
	t.Run("SingleTest", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().DeleteTest("00000000-0000-0000-0000-000000000000").Return(nil).Times(1)

		loadtestTestDelete(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MultipleTests", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteTest("00000000-0000-0000-0000-000000000000").Return(nil),
			service.EXPECT().DeleteTest("00000000-0000-0000-0000-000000000001").Return(nil),
		)

		loadtestTestDelete(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("GetTestError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().DeleteTest("00000000-0000-0000-0000-000000000000").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing test [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			loadtestTestDelete(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}
