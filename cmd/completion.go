package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func CompletionRootCmd() *cobra.Command {
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
		Long: `To load completion into current shell:

source <(ans completion bash)

To configure your bash shell to load completions for all sessions, add the above to ~/.bashrc as below:

echo 'source <(ans completion bash)' >> ~/.bashrc

Alternatively, completions in /etc/bash_completion.d/ will be auto-loaded:

echo 'source <(ans completion bash)' >> /etc/bash_completion.d/ans
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
		Long: `To load completion into current shell:

Invoke-Expression -Command (ans completion powershell | Out-String)

To configure your shell to load completions for each session, output completion to profile:

$ProfileDIR = Split-Path -Parent -Path $PROFILE
$CompletionPath = [System.IO.Path]::GetFullPath("$ProfileDIR/ans.completion.ps1")
Out-File -Append -FilePath $CompletionPath -Encoding ASCII -InputObject "Invoke-Expression -Command (ans completion powershell | Out-String)"
Out-File -Append -FilePath $PROFILE -Encoding ASCII -InputObject ` + "\"`n. $CompletionPath\"",
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenPowerShellCompletion(os.Stdout)
		},
	}
}

func completionZshCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "zsh",
		Short: "Generates zsh completion scripts",
		Long: `To load completion into current shell:

source <(ans completion zsh)

To configure your zsh shell to load completions for all sessions, completions in /etc/bash_completion.d/ will be auto-loaded:

echo 'source <(ans completion zsh)' >> /etc/bash_completion.d/ans
`,
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenZshCompletion(os.Stdout)
		},
	}
}
