package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func completionRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Commands for generating shell completions",
	}

	// Child commands
	cmd.AddCommand(completionBashCmd())
	cmd.AddCommand(completionPowerShellCmd())
	cmd.AddCommand(completionZshCmd())

	return cmd
}

func completionBashCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bash",
		Short: "Generates bash completion scripts",
		Long: `To load completion run

. <(ukfast completion bash)

To configure your bash shell to load completions for each session, output completion to file:

ukfast completion bash > /etc/bash_completion.d/ukfast
`,
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenBashCompletion(os.Stdout)
		},
	}
}

func completionPowerShellCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "powershell",
		Short: "Generates Powershell completion scripts",
		Long: `To load completion run

Invoke-Expression -command (ukfast completion powershell | Out-String)

To configure your shell to load completions for each session, output completion to profile:

$ProfileDIR = Split-Path -Parent -Path $profile
$CompletionPath = [System.IO.Path]::GetFullPath("$ProfileDIR/ukfast.completion.ps1")
ukfast completion powershell | Out-File -Append -FilePath $CompletionPath -Encoding ASCII
` + "\"`n. $CompletionPath\" | Out-File -Append -FilePath $profile -Encoding ASCII",
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenPowerShellCompletion(os.Stdout)
		},
	}
}

func completionZshCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "zsh",
		Short: "Generates zsh completion scripts",
		Long: `To load completion run

. <(ukfast completion zsh)

To configure your zsh shell to load completions for each session, output completion to file:

ukfast completion zsh > /etc/bash_completion.d/ukfast
`,
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenZshCompletion(os.Stdout)
		},
	}
}
