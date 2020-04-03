package safedns

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func Test_safednsZoneNoteListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := safednsZoneNoteListCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.com"})

		assert.Nil(t, err)
	})

	t.Run("MissingZone_Error", func(t *testing.T) {
		cmd := safednsZoneNoteListCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})
}

func Test_safednsZoneNoteList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetZoneNotes("testdomain1.com", gomock.Any()).Return([]safedns.Note{}, nil).Times(1)

		safednsZoneNoteList(service, &cobra.Command{}, []string{"testdomain1.com"})
	})

	t.Run("ExpectedFilterFromFlags", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsZoneNoteListCmd(nil)
		cmd.Flags().Set("ip", "1.2.3.4")

		expectedParams := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				connection.APIRequestFiltering{
					Property: "ip",
					Operator: connection.EQOperator,
					Value:    []string{"1.2.3.4"},
				},
			},
		}

		service.EXPECT().GetZoneNotes("testdomain1.com", gomock.Eq(expectedParams)).Return([]safedns.Note{}, nil).Times(1)

		safednsZoneNoteList(service, cmd, []string{"testdomain1.com"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := safednsZoneNoteList(service, cmd, []string{"123"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetZonesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetZoneNotes("testdomain1.com", gomock.Any()).Return([]safedns.Note{}, errors.New("test error")).Times(1)

		err := safednsZoneNoteList(service, &cobra.Command{}, []string{"testdomain1.com"})

		assert.Equal(t, "Error retrieving notes for zone: test error", err.Error())
	})
}

func Test_safednsZoneNoteShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := safednsZoneNoteShowCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.com", "123"})

		assert.Nil(t, err)
	})

	t.Run("MissingZone_Error", func(t *testing.T) {
		cmd := safednsZoneNoteShowCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})

	t.Run("MissingNote_Error", func(t *testing.T) {
		cmd := safednsZoneNoteShowCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.com"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing note", err.Error())
	})
}

func Test_safednsZoneNoteShow(t *testing.T) {
	t.Run("SingleZoneNote", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetZoneNote("testdomain1.com", 123).Return(safedns.Note{}, nil).Times(1)

		safednsZoneNoteShow(service, &cobra.Command{}, []string{"testdomain1.com", "123"})
	})

	t.Run("MultipleZoneNotes", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetZoneNote("testdomain1.com", 123).Return(safedns.Note{}, nil),
			service.EXPECT().GetZoneNote("testdomain1.com", 456).Return(safedns.Note{}, nil),
		)

		safednsZoneNoteShow(service, &cobra.Command{}, []string{"testdomain1.com", "123", "456"})
	})

	t.Run("InvalidNoteID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid note ID [abc]\n", func() {
			safednsZoneNoteShow(service, &cobra.Command{}, []string{"testdomain1.com", "abc"})
		})
	})

	t.Run("GetZoneNoteError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().GetZoneNote("testdomain1.com", 123).Return(safedns.Note{}, errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error retrieving note [123]: test error\n", func() {
			safednsZoneNoteShow(service, &cobra.Command{}, []string{"testdomain1.com", "123"})
		})
	})
}

func Test_safednsZoneNoteCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := safednsZoneNoteCreateCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.com"})

		assert.Nil(t, err)
	})

	t.Run("MissingZone_Error", func(t *testing.T) {
		cmd := safednsZoneNoteCreateCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing zone", err.Error())
	})
}

func Test_safednsZoneNoteCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)
		cmd := safednsZoneNoteCreateCmd(nil)
		cmd.Flags().Set("notes", "test note 1")
		cmd.Flags().Set("ip", "1.2.3.4")

		expectedRequest := safedns.CreateNoteRequest{
			Notes: "test note 1",
		}

		gomock.InOrder(
			service.EXPECT().CreateZoneNote("testdomain1.com", expectedRequest).Return(123, nil),
			service.EXPECT().GetZoneNote("testdomain1.com", 123).Return(safedns.Note{}, nil),
		)

		safednsZoneNoteCreate(service, cmd, []string{"testdomain1.com"})
	})

	t.Run("CreateZoneNoteError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		service.EXPECT().CreateZoneNote("testdomain1.com", gomock.Any()).Return(0, errors.New("test error")).Times(1)

		err := safednsZoneNoteCreate(service, &cobra.Command{}, []string{"testdomain1.com"})

		assert.Equal(t, "Error creating note: test error", err.Error())
	})

	t.Run("GetZoneNoteError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSafeDNSService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().CreateZoneNote("testdomain1.com", gomock.Any()).Return(123, nil),
			service.EXPECT().GetZoneNote("testdomain1.com", 123).Return(safedns.Note{}, errors.New("test error")),
		)

		err := safednsZoneNoteCreate(service, &cobra.Command{}, []string{"testdomain1.com"})

		assert.Equal(t, "Error retrieving new note: test error", err.Error())
	})
}
