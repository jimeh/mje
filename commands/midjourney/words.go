package midjourney

import (
	"sort"

	"github.com/jimeh/go-midjourney"
	"github.com/jimeh/mje/commands/render"
	"github.com/jimeh/mje/commands/shared"
	"github.com/spf13/cobra"
)

type Word struct {
	Word     string `json:"word"`
	ImageURL string `json:"image_url"`
}

func NewWords(mc *midjourney.Client) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "words",
		Aliases: []string{"w"},
		Short:   "Get dictionary words",
		RunE:    wordsRunE(mc),
	}

	cmd.Flags().StringP("format", "f", "", "output format (yaml or json)")
	cmd.Flags().StringP("query", "q", "", "query to search for")
	cmd.Flags().IntP("amount", "a", 50, "amount of words to fetch")
	cmd.Flags().IntP("page", "p", 0, "page to fetch")
	cmd.Flags().IntP("seed", "s", 0, "seed")

	return cmd, nil
}

func wordsRunE(mc *midjourney.Client) shared.RunEFunc {
	return func(cmd *cobra.Command, _ []string) error {
		fs := cmd.Flags()
		q := &midjourney.WordsQuery{}

		if v, err := fs.GetString("query"); err == nil && v != "" {
			q.Query = v
		}
		if v, err := fs.GetInt("amount"); err == nil && v > 0 {
			q.Amount = v
		}
		if v, err := fs.GetInt("page"); err == nil && v != 0 {
			q.Page = v
		}
		if v, err := fs.GetInt("seed"); err == nil && v != 0 {
			q.Seed = v
		}

		words, err := mc.Words(cmd.Context(), q)
		if err != nil {
			return err
		}

		r := []*Word{}
		for _, w := range words {
			r = append(r, &Word{
				Word:     w.Word,
				ImageURL: w.ImageURL(),
			})
		}

		format := shared.FlagString(cmd, "format")
		sort.SliceStable(r, func(i, j int) bool {
			return r[i].Word < r[j].Word
		})

		return render.Render(cmd.OutOrStdout(), format, r)
	}
}
