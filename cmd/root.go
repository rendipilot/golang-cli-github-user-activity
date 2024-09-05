package cmd

import (
	"fmt"
	"pilotkode/github-user-activity/activity"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "github-activity",
		Short: "Github user activity adalah alat cli untuk mengambil data aktivitas dari sebuah user github",
		Long: `Github user activity adalah alat cli untuk mengambil data aktivitas dari sebuah user github, Kita dapat melihat
		       data aktivitas seperti push atau pull yang dilakukan`,
		RunE: func (cmd *cobra.Command, args []string) error {
			// todo return display
			return RunDisplayActivityCmd(args)
		},
	}

	return cmd;
}

func RunDisplayActivityCmd(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("please provide a username")
	}

	username := args[0]
	
	// get activity user
	activities, code, err := activity.FetchGithubActivity(username)
	
	if err != nil {
        return fmt.Errorf("error fetching user activity: %v", err)
    }

	return activity.DisplayActivity(username, activities, code)

}

