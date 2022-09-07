package commands

import (
	"io"
	"os"

	"github.com/jimeh/mj2n/midjourney"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type runEFunc func(cmd *cobra.Command, _ []string) error

func NewMJ2N() (*cobra.Command, error) {
	mc, err := midjourney.New(midjourney.WithUserAgent("mj2n/0.0.1-dev"))
	if err != nil {
		return nil, err
	}

	cmd := &cobra.Command{
		Use:               "mj2n",
		Short:             "MidJourney to Notion importer",
		PersistentPreRunE: persistentPreRunE(mc),
	}

	cmd.PersistentFlags().StringP(
		"log-level", "l", "info",
		"one of: trace, debug, info, warn, error, fatal, panic",
	)
	cmd.PersistentFlags().StringP(
		"mj-token", "m", "", "MidJourney API token",
	)
	cmd.PersistentFlags().String(
		"mj-api-url", midjourney.DefaultAPIURL.String(), "MidJourney API URL",
	)

	midjourneyCmd, err := NewMidjourney(mc)
	if err != nil {
		return nil, err
	}

	cmd.AddCommand(midjourneyCmd)

	return cmd, nil
}

func persistentPreRunE(mc *midjourney.Client) runEFunc {
	return func(cmd *cobra.Command, _ []string) error {
		err := setupZerolog(cmd)
		if err != nil {
			return err
		}

		err = setupMidJourney(cmd, mc)
		if err != nil {
			return err
		}

		return nil
	}
}

func setupMidJourney(cmd *cobra.Command, mc *midjourney.Client) error {
	opts := []midjourney.Option{
		midjourney.WithLogger(log.Logger),
	}

	if f := cmd.Flag("mj-token"); f.Changed {
		opts = append(opts, midjourney.WithAuthToken(f.Value.String()))
	} else if v := os.Getenv("MIDJOURNEY_TOKEN"); v != "" {
		opts = append(opts, midjourney.WithAuthToken(v))
	}

	apiURL := flagString(cmd, "mj-api-url")
	if apiURL == "" {
		apiURL = os.Getenv("MIDJOURNEY_API_URL")
	}
	if apiURL != "" {
		opts = append(opts, midjourney.WithAPIURL(apiURL))
	}

	return mc.Set(opts...)
}

func setupZerolog(cmd *cobra.Command) error {
	var levelStr string
	if v := os.Getenv("MJ2N_DEBUG"); v != "" {
		levelStr = "debug"
	} else if v := os.Getenv("MJ2N_LOG_LEVEL"); v != "" {
		levelStr = v
	}

	var out io.Writer = os.Stderr

	if cmd != nil {
		out = cmd.OutOrStderr()
		fl := cmd.Flag("log-level")
		if fl != nil && (fl.Changed || levelStr == "") {
			levelStr = fl.Value.String()
		}
	}

	if levelStr == "" {
		levelStr = "info"
	}

	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		return err
	}

	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = ""

	output := zerolog.ConsoleWriter{Out: out}
	output.FormatTimestamp = func(i interface{}) string {
		return ""
	}

	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	return nil
}

func flagString(cmd *cobra.Command, name string) string {
	var r string

	if f := cmd.Flag(name); f != nil {
		r = f.Value.String()
	}

	return r
}
