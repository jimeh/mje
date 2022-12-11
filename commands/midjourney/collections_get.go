package midjourney

import (
	"github.com/jimeh/go-midjourney"
	"github.com/jimeh/mje/commands/render"
	"github.com/jimeh/mje/commands/shared"
	"github.com/spf13/cobra"
)

func NewCollectionsGet(mc *midjourney.Client) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:        "get collection_id",
		Short:      "Get a collection",
		RunE:       collectionsGetRunE(mc),
		ArgAliases: []string{"collection_id"},
		Args:       cobra.ExactArgs(1),
	}

	return cmd, nil
}

func collectionsGetRunE(mc *midjourney.Client) shared.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		q := &midjourney.CollectionsQuery{
			CollectionID: args[0],
		}

		cols, err := mc.Collections(cmd.Context(), q)
		if err != nil {
			return err
		}

		format := shared.FlagString(cmd, "format")

		return render.Render(cmd.OutOrStdout(), format, cols)
	}
}
