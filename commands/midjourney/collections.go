package midjourney

import (
	"github.com/jimeh/go-midjourney"
	"github.com/spf13/cobra"
)

func NewCollections(mc *midjourney.Client) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "collections",
		Aliases: []string{"collection", "col"},
		Short:   "Query collections",
	}

	listCmd, err := NewCollectionsList(mc)
	if err != nil {
		return nil, err
	}
	getCmd, err := NewCollectionsGet(mc)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(
		listCmd,
		getCmd,
	)

	return cmd, nil
}
