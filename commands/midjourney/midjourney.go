package midjourney

import (
	"github.com/jimeh/go-midjourney"
	"github.com/spf13/cobra"
)

func New(mc *midjourney.Client) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "midjourney",
		Aliases: []string{"mj"},
		Short:   "Query the Midjourney API directly",
	}

	collectionsCmd, err := NewCollections(mc)
	if err != nil {
		return nil, err
	}
	recentJobsCmd, err := NewRecentJobs(mc)
	if err != nil {
		return nil, err
	}
	wordsCmd, err := NewWords(mc)
	if err != nil {
		return nil, err
	}

	cmd.AddCommand(
		collectionsCmd,
		recentJobsCmd,
		wordsCmd,
	)

	return cmd, nil
}
