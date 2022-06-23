package cmd

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

func updateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Updates CLI to latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			currentVersion, err := semver.ParseTolerant(appVersion)
			if err != nil {
				return fmt.Errorf("Unable to parse version: %s", err.Error())
			}
			newRelease, err := selfupdate.UpdateSelf(currentVersion, "ans-group/cli")
			if err != nil {
				return fmt.Errorf("Failed to update ANS CLI: %s", err.Error())
			}

			if currentVersion.Equals(newRelease.Version) {
				fmt.Printf("ANS CLI already at latest version (%s)\n", appVersion)
			} else {
				fmt.Printf("ANS CLI updated to version v%s successfully\n", newRelease.Version)
				fmt.Println("Release notes:\n", newRelease.ReleaseNotes)
			}
			return nil
		},
	}
}
