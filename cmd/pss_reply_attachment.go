package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/pss"
)

func pssReplyAttachmentRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attachment",
		Short: "sub-commands relating to reply attachments",
	}

	// Child commands
	cmd.AddCommand(pssReplyAttachmentDownloadCmd())

	return cmd
}

func pssReplyAttachmentDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "download <reply: id> <attachment: name>",
		Short:   "Shows a reply",
		Long:    "This command shows one or more replies",
		Example: "ukfast pss reply attachment download 123 file.txt --path /my/file.txt",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing reply")
			}
			if len(args) < 2 {
				return errors.New("Missing attachment")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			pssReplyAttachmentDownload(getClient().PSSService(), cmd, args)
		},
	}

	cmd.Flags().String("path", "", "Specifies destination path for file. Omitting this flag will save in current working directory")

	return cmd
}

func pssReplyAttachmentDownload(service pss.PSSService, cmd *cobra.Command, args []string) {
	attachmentStream, err := service.DownloadReplyAttachmentStream(args[0], args[1])
	if err != nil {
		output.Fatalf("Error downloading reply attachment: %s", err)
		return
	}

	var targetFilePath string
	if cmd.Flags().Changed("path") {
		targetFilePath, _ := cmd.Flags().GetString("path")
		targetFilePath = filepath.Join(targetFilePath, args[1])
	} else {
		dir, err := os.Getwd()
		if err != nil {
			output.Fatalf("Error determining current directory: %s", err)
			return
		}
		targetFilePath = filepath.Join(dir, args[1])
	}

	_, err = appFilesystem.Stat(targetFilePath)
	if err == nil || !os.IsNotExist(err) {
		output.Fatalf("Destination file [%s] exists", targetFilePath)
		return
	}

	err = afero.WriteReader(appFilesystem, targetFilePath, attachmentStream)
	if err != nil {
		output.Fatalf("Error writing attachment to [%s]: %s", targetFilePath, err.Error())
		return
	}
}
