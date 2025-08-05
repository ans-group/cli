package safedns

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/safedns"
	"github.com/spf13/cobra"
)

func safednsZoneNoteRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "note",
		Short: "sub-commands relating to zone notes",
	}

	// Child commands
	cmd.AddCommand(safednsZoneNoteListCmd(f))
	cmd.AddCommand(safednsZoneNoteShowCmd(f))
	cmd.AddCommand(safednsZoneNoteCreateCmd(f))

	return cmd
}

func safednsZoneNoteListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <zone: name>",
		Short:   "Lists zone notes",
		Long:    "This command lists zone notes",
		Example: "ans safedns zone note list ans.co.uk\nans safedns zone note list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing zone")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return safednsZoneNoteList(c.SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().String("ip", "", "Zone note IP address for filtering")

	return cmd
}

func safednsZoneNoteList(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("ip", "ip"))
	if err != nil {
		return err
	}

	zoneNotes, err := service.GetZoneNotes(args[0], params)
	if err != nil {
		return fmt.Errorf("error retrieving notes for zone: %s", err)
	}

	return output.CommandOutput(cmd, NoteCollection(zoneNotes))
}

func safednsZoneNoteShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <zone: name> <note: id>...",
		Short:   "Shows a zone note",
		Long:    "This command shows one or more zone notes",
		Example: "ans safedns zone note show ans.co.uk 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing zone")
			}
			if len(args) < 2 {
				return errors.New("missing note")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return safednsZoneNoteShow(c.SafeDNSService(), cmd, args)
		},
	}
}

func safednsZoneNoteShow(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	var zoneNotes []safedns.Note

	for _, arg := range args[1:] {
		noteID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid note ID [%s]", arg)
			continue
		}

		zoneNote, err := service.GetZoneNote(args[0], noteID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving note [%d]: %s", noteID, err)
			continue
		}

		zoneNotes = append(zoneNotes, zoneNote)
	}

	return output.CommandOutput(cmd, NoteCollection(zoneNotes))
}

func safednsZoneNoteCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <zone: name>",
		Short:   "Creates a zone note",
		Long:    "This command creates a zone note",
		Example: "ans safedns zone note create ans.co.uk --notes \"test note\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing zone")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return safednsZoneNoteCreate(c.SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().Int("contact-id", 0, "Contact ID for note")
	cmd.Flags().String("notes", "", "Note content")
	_ = cmd.MarkFlagRequired("notes")

	return cmd
}

func safednsZoneNoteCreate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	createRequest := safedns.CreateNoteRequest{}
	createRequest.ContactID, _ = cmd.Flags().GetInt("contact-id")
	createRequest.Notes, _ = cmd.Flags().GetString("notes")

	id, err := service.CreateZoneNote(args[0], createRequest)
	if err != nil {
		return fmt.Errorf("error creating note: %s", err)
	}

	zoneNote, err := service.GetZoneNote(args[0], id)
	if err != nil {
		return fmt.Errorf("error retrieving new note: %s", err)
	}

	return output.CommandOutput(cmd, NoteCollection([]safedns.Note{zoneNote}))
}
