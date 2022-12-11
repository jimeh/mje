package shared

import "github.com/spf13/cobra"

func FlagString(cmd *cobra.Command, name string) string {
	var r string

	if f := cmd.Flag(name); f != nil {
		r = f.Value.String()
	}

	return r
}
