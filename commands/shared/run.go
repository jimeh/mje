package shared

import "github.com/spf13/cobra"

type RunEFunc func(cmd *cobra.Command, _ []string) error
