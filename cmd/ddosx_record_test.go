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
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func Test_ddosxRecordList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetRecords(gomock.Any()).Return([]ddosx.Record{}, nil).Times(1)

		ddosxRecordList(service, &cobra.Command{}, []string{})
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

		service := mocks.NewMockDDoSXService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		output := test.CatchStdErr(t, func() {
			ddosxRecordList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Missing value for filtering\n", output)
	})

	t.Run("GetRecordsError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetRecords(gomock.Any()).Return([]ddosx.Record{}, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			ddosxRecordList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving records: test error\n", output)
	})
}
