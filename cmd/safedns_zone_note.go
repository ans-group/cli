package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func safednsZoneNoteRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "note",
		Short: "sub-commands relating to zone notes",
	}

	// Child commands
	cmd.AddCommand(safednsZoneNoteListCmd())
	cmd.AddCommand(safednsZoneNoteShowCmd())
	cmd.AddCommand(safednsZoneNoteCreateCmd())

	return cmd
}

func safednsZoneNoteListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <zone: name>",
		Short:   "Lists zone notes",
		Long:    "This command lists zone notes",
		Example: "ukfast safedns zone note list ukfast.co.uk\nukfast safedns zone note list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing zone")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			safednsZoneNoteList(getClient().SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().String("ip", "", "Zone note IP address for filtering")

	return cmd
}

func safednsZoneNoteList(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	if cmd.Flags().Changed("ip") {
		filterIP, _ := cmd.Flags().GetString("ip")
		params.WithFilter(helper.GetFilteringInferOperator("ip", filterIP))
	}

	zoneNotes, err := service.GetZoneNotes(args[0], params)
	if err != nil {
		output.Fatalf("Error retrieving notes for zone: %s", err)
		return
	}

	outputSafeDNSNotes(zoneNotes)
}

func safednsZoneNoteShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <zone: name> <note: id>...",
		Short:   "Shows a zone note",
		Long:    "This command shows one or more zone notes",
		Example: "ukfast safedns zone note show ukfast.co.uk 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing zone")
			}
			if len(args) < 2 {
				return errors.New("Missing note")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			safednsZoneNoteShow(getClient().SafeDNSService(), cmd, args)
		},
	}
}

func safednsZoneNoteShow(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
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

	outputSafeDNSNotes(zoneNotes)
}

func safednsZoneNoteCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <zone: name>",
		Short:   "Creates a zone note",
		Long:    "This command creates a zone note",
		Example: "ukfast safedns zone note create ukfast.co.uk --notes \"test note\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing zone")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			safednsZoneNoteCreate(getClient().SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().Int("contact-id", 0, "Contact ID for note")
	cmd.Flags().String("notes", "", "Note content")
	cmd.MarkFlagRequired("notes")

	return cmd
}

func safednsZoneNoteCreate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	createRequest := safedns.CreateNoteRequest{}
	createRequest.ContactID, _ = cmd.Flags().GetInt("contact-id")
	createRequest.Notes, _ = cmd.Flags().GetString("notes")

	id, err := service.CreateZoneNote(args[0], createRequest)
	if err != nil {
		output.Fatalf("Error creating note: %s", err)
		return
	}

	zoneNote, err := service.GetZoneNote(args[0], id)
	if err != nil {
		output.Fatalf("Error retrieving new note: %s", err)
		return
	}

	outputSafeDNSNotes([]safedns.Note{zoneNote})
}
