package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/ukfast/cli/internal/pkg/output"

	"github.com/spf13/cobra"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainVerificationFileUploadRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fileupload",
		Short: "sub-commands relating to file upload domain verification",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainVerificationFileUploadShowCmd())
	cmd.AddCommand(ddosxDomainVerificationFileUploadDownloadCmd())

	return cmd
}

func ddosxDomainVerificationFileUploadShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name>...",
		Short:   "Shows the verification file for a domain",
		Long:    "This command shows the verification file for one or more domains, for use with the file upload verification method",
		Example: "ukfast ddosx domain verification fileupload show example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainVerificationFileUploadShow(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainVerificationFileUploadShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var files []OutputDDoSXDomainVerificationFilesFile

	for _, arg := range args {
		content, filename, err := service.DownloadDomainVerificationFile(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving domain verification file [%s]: %s", arg, err)
			continue
		}

		files = append(files, OutputDDoSXDomainVerificationFilesFile{
			Name:    filename,
			Content: content,
		})
	}

	outputDDoSXDomainVerificationFiles(files)
}

func ddosxDomainVerificationFileUploadDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "download <domain: name>",
		Short:   "Downloads the verification file for a domain",
		Long:    "This command downloads the verification file for a domain, for use with the file upload verification method",
		Example: "ukfast ddosx domain verification fileupload download example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainVerificationFileUploadDownload(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("path", "", "Target directory path for file download. File name is automatically determined by command")
	cmd.MarkFlagRequired("path")

	return cmd
}

func ddosxDomainVerificationFileUploadDownload(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	content, filename, err := service.DownloadDomainVerificationFile(args[0])
	if err != nil {
		output.Fatalf("Error retrieving domain verification file: %s", err)
		return
	}

	directory, _ := cmd.Flags().GetString("path")

	targetFilePath := filepath.Join(directory, filename)

	_, err = appFilesystem.Stat(targetFilePath)
	if err == nil || !os.IsNotExist(err) {
		output.Fatalf("Destination file [%s] exists", targetFilePath)
		return
	}

	err = afero.WriteFile(appFilesystem, targetFilePath, []byte(content), 0644)
	if err != nil {
		output.Fatalf("Error writing domain verification file to [%s]: %s", targetFilePath, err.Error())
		return
	}

	fmt.Println(targetFilePath)
}
