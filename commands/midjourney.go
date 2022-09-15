package commands

import (
	"github.com/jimeh/go-midjourney"
	"github.com/spf13/cobra"
)

func NewMidjourney(mc *midjourney.Client) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "midjourney",
		Aliases: []string{"mj"},
		Short:   "Query the Midjourney API directly",
	}

	recentJobsCmd, err := NewMidjourneyRecentJobs(mc)
	if err != nil {
		return nil, err
	}
	wordsCmd, err := NewMidjourneyWords(mc)
	if err != nil {
		return nil, err
	}

	cmd.AddCommand(
		recentJobsCmd,
		wordsCmd,
	)

	return cmd, nil
}
