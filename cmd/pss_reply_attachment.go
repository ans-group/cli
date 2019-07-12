package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
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
	cmd.AddCommand(pssReplyAttachmentUploadCmd())

	return cmd
}

func pssReplyAttachmentDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "download <reply: id> <attachment: name>",
		Short:   "Downloads a reply attachment",
		Long:    "This command downloads a reply attachment",
		Example: "ukfast pss reply attachment download 123 file.txt --path /path/to/file",
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

	path, _ := cmd.Flags().GetString("path")
	targetFilePath, err := helper.GetDestinationFilePath(appFilesystem, args[1], path)
	if err != nil {
		output.Fatalf("Error determining destination file path: %s", err)
		return
	}

	_, err = appFilesystem.Stat(targetFilePath)
	if err == nil || !os.IsNotExist(err) {
		output.Fatalf("Destination file [%s] exists", targetFilePath)
		return
	}

	err = afero.SafeWriteReader(appFilesystem, targetFilePath, attachmentStream)
	if err != nil {
		output.Fatalf("Error writing attachment to [%s]: %s", targetFilePath, err.Error())
		return
	}
}

func pssReplyAttachmentUploadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upload <reply: id>",
		Short:   "Uploads a reply attachment",
		Long:    "This command uploads a reply attachment",
		Example: "ukfast pss reply attachment upload 123 --path /path/to/file",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing reply")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := pssReplyAttachmentUpload(getClient().PSSService(), cmd, args)
			if err != nil {
				output.Fatal(err.Error())
			}
		},
	}

	cmd.Flags().String("path", "", "Specifies path for file to upload")
	cmd.MarkFlagRequired("path")

	return cmd
}

func pssReplyAttachmentUpload(service pss.PSSService, cmd *cobra.Command, args []string) error {
	path, _ := cmd.Flags().GetString("path")

	fileStream, err := appFilesystem.Open(path)
	if err != nil {
		return fmt.Errorf("Failed to open file [%s]: %s", path, err)
	}

	err = service.UploadReplyAttachmentStream(args[0], filepath.Base(path), fileStream)
	if err != nil {
		return fmt.Errorf("Failed to upload attachment: %s", err)
	}

	return nil
}
