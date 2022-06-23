package ddosx

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"

	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
)

func ddosxDomainVerificationFileUploadRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fileupload",
		Short: "sub-commands relating to file upload domain verification",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainVerificationFileUploadShowCmd(f))
	cmd.AddCommand(ddosxDomainVerificationFileUploadDownloadCmd(f, fs))
	cmd.AddCommand(ddosxDomainVerificationFileUploadVerifyCmd(f))

	return cmd
}

func ddosxDomainVerificationFileUploadShowCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainVerificationFileUploadShow(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainVerificationFileUploadShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var files []OutputDDoSXDomainVerificationFilesFile

	for _, arg := range args {
		content, filename, err := service.DownloadDomainVerificationFile(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain verification file [%s]: %s", arg, err)
			continue
		}

		files = append(files, OutputDDoSXDomainVerificationFilesFile{
			Name:    filename,
			Content: content,
		})
	}

	return output.CommandOutput(cmd, OutputDDoSXDomainVerificationFilesProvider(files))
}

func ddosxDomainVerificationFileUploadDownloadCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainVerificationFileUploadDownload(c.DDoSXService(), fs, cmd, args)
		},
	}

	cmd.Flags().String("path", "", "Target directory path for file download. File name is automatically determined by command")
	cmd.MarkFlagRequired("path")

	return cmd
}

func ddosxDomainVerificationFileUploadDownload(service ddosx.DDoSXService, fs afero.Fs, cmd *cobra.Command, args []string) error {
	content, filename, err := service.DownloadDomainVerificationFile(args[0])
	if err != nil {
		return fmt.Errorf("Error retrieving domain verification file: %s", err)
	}

	directory, _ := cmd.Flags().GetString("path")

	targetFilePath := filepath.Join(directory, filename)

	_, err = fs.Stat(targetFilePath)
	if err == nil || !os.IsNotExist(err) {
		return fmt.Errorf("Destination file [%s] exists", targetFilePath)
	}

	err = afero.WriteFile(fs, targetFilePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("Error writing domain verification file to [%s]: %s", targetFilePath, err.Error())
	}

	fmt.Println(targetFilePath)
	return nil
}

func ddosxDomainVerificationFileUploadVerifyCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "verify <domain: name>...",
		Short:   "Verifies a domain via verification file method",
		Long:    "This command verifies one or more domains via the verification file method",
		Example: "ukfast ddosx domain verification fileupload verify example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainVerificationFileUploadVerify(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainVerificationFileUploadVerify(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.VerifyDomainFileUpload(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error verifying domain [%s] via verification file method: %s", arg, err)
			continue
		}
	}

	return nil
}
