package commands

import (
	"github.com/jimeh/mj2n/midjourney"
	"github.com/spf13/cobra"
)

func NewMidjourney(mc *midjourney.Client) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "midjourney",
		Aliases: []string{"mj"},
		Short:   "MidJourney specific commands",
	}

	recentJobsCmd, err := NewMidjourneyRecentJobs(mc)
	if err != nil {
		return nil, err
	}

	cmd.AddCommand(recentJobsCmd)

	return cmd, nil
}
