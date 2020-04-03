package pss

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/pss"
)

func pssReplyAttachmentRootCmd(f factory.ClientFactory, appFilesystem afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attachment",
		Short: "sub-commands relating to reply attachments",
	}

	// Child commands
	cmd.AddCommand(pssReplyAttachmentDownloadCmd(f, appFilesystem))
	cmd.AddCommand(pssReplyAttachmentUploadCmd(f, appFilesystem))
	cmd.AddCommand(pssReplyAttachmentDeleteCmd(f))

	return cmd
}

func pssReplyAttachmentDownloadCmd(f factory.ClientFactory, appFilesystem afero.Fs) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return pssReplyAttachmentDownload(f.NewClient().PSSService(), appFilesystem, cmd, args)
		},
	}

	cmd.Flags().String("path", "", "Specifies destination path for file. Omitting this flag will save in current working directory")

	return cmd
}

func pssReplyAttachmentDownload(service pss.PSSService, appFilesystem afero.Fs, cmd *cobra.Command, args []string) error {
	attachmentStream, err := service.DownloadReplyAttachmentStream(args[0], args[1])
	if err != nil {
		return fmt.Errorf("Error downloading reply attachment: %s", err)
	}

	path, _ := cmd.Flags().GetString("path")
	targetFilePath, err := helper.GetDestinationFilePath(appFilesystem, args[1], path)
	if err != nil {
		return fmt.Errorf("Error determining destination file path: %s", err)
	}

	_, err = appFilesystem.Stat(targetFilePath)
	if err == nil || !os.IsNotExist(err) {
		return fmt.Errorf("Destination file [%s] exists", targetFilePath)
	}

	err = afero.SafeWriteReader(appFilesystem, targetFilePath, attachmentStream)
	if err != nil {
		return fmt.Errorf("Error writing attachment to [%s]: %s", targetFilePath, err.Error())
	}

	return nil
}

func pssReplyAttachmentUploadCmd(f factory.ClientFactory, appFilesystem afero.Fs) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return pssReplyAttachmentUpload(f.NewClient().PSSService(), appFilesystem, cmd, args)
		},
	}

	cmd.Flags().String("path", "", "Specifies path for file to upload")
	cmd.MarkFlagRequired("path")

	return cmd
}

func pssReplyAttachmentUpload(service pss.PSSService, appFilesystem afero.Fs, cmd *cobra.Command, args []string) error {
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

func pssReplyAttachmentDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <reply: id> <attachment: name>...",
		Short:   "Deletes a reply attachment",
		Long:    "This command deletes one or more reply attachments",
		Example: "ukfast pss reply attachment delete 123 file.txt",
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
			pssReplyAttachmentDelete(f.NewClient().PSSService(), cmd, args)
		},
	}
}

func pssReplyAttachmentDelete(service pss.PSSService, cmd *cobra.Command, args []string) {
	for _, arg := range args[1:] {
		err := service.DeleteReplyAttachment(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error deleting reply attachment [%s]: %s", arg, err)
		}
	}
}
