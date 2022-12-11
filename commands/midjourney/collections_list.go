package midjourney

import (
	"github.com/jimeh/go-midjourney"
	"github.com/jimeh/mje/commands/render"
	"github.com/jimeh/mje/commands/shared"
	"github.com/spf13/cobra"
)

func NewCollectionsList(mc *midjourney.Client) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List collections",
		RunE:    collectionsListRunE(mc),
	}

	cmd.Flags().StringP("user-id", "u", "", "user ID to list jobs for")
	cmd.Flags().StringP(
		"collection-id", "c", "", "collection ID to list jobs for",
	)

	return cmd, nil
}

func collectionsListRunE(mc *midjourney.Client) shared.RunEFunc {
	return func(cmd *cobra.Command, _ []string) error {
		fs := cmd.Flags()
		q := &midjourney.CollectionsQuery{}

		if v, err := fs.GetString("user-id"); err == nil && v != "" {
			q.UserID = v
		}
		if v, err := fs.GetString("collection-id"); err == nil && v != "" {
			q.CollectionID = v
		}
		cols, err := mc.Collections(cmd.Context(), q)
		if err != nil {
			return err
		}

		format := shared.FlagString(cmd, "format")

		return render.Render(cmd.OutOrStdout(), format, cols)
	}
}
