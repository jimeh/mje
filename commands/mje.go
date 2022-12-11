package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/jimeh/go-midjourney"
	mjcmds "github.com/jimeh/mje/commands/midjourney"
	"github.com/jimeh/mje/commands/shared"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type Info struct {
	Version string
	Commit  string
	Date    string
}

func New(info Info) (*cobra.Command, error) {
	if info.Version == "" {
		info.Version = "0.0.0-dev"
	}

	mc, err := midjourney.New(midjourney.WithUserAgent("mje/" + info.Version))
	if err != nil {
		return nil, err
	}

	cmd := &cobra.Command{
		Use:               "mje",
		Short:             "MidJourney exporter",
		Version:           info.Version,
		PersistentPreRunE: persistentPreRunE(mc),
	}

	cmd.PersistentFlags().String(
		"log-level", "info",
		"one of: trace, debug, info, warn, error, fatal, panic",
	)
	cmd.PersistentFlags().String(
		"log-format", "plain",
		"one of: plain, json",
	)
	cmd.PersistentFlags().String(
		"token", "", "MidJourney token",
	)
	cmd.PersistentFlags().String(
		"api-url", midjourney.DefaultAPIURL.String(),
		"MidJourney API URL",
	)

	midjourneyCmd, err := mjcmds.New(mc)
	if err != nil {
		return nil, err
	}

	cmd.AddCommand(midjourneyCmd)

	return cmd, nil
}

func persistentPreRunE(mc *midjourney.Client) shared.RunEFunc {
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

	if f := cmd.Flag("token"); f != nil && f.Changed {
		opts = append(opts, midjourney.WithAuthToken(f.Value.String()))
	} else if v := os.Getenv("MIDJOURNEY_TOKEN"); v != "" {
		opts = append(opts, midjourney.WithAuthToken(v))
	}

	apiURL := shared.FlagString(cmd, "api-url")
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
	var logFormat string

	if v := os.Getenv("MJE_DEBUG"); v != "" {
		levelStr = "debug"
	} else if v := os.Getenv("MJE_LOG_LEVEL"); v != "" {
		levelStr = v
	}
	if v := os.Getenv("MJE_LOG_FORMAT"); v != "" {
		logFormat = v
	}

	var out io.Writer = os.Stderr

	if cmd != nil {
		out = cmd.OutOrStderr()
		fl := cmd.Flag("log-level")
		if fl != nil && (fl.Changed || levelStr == "") {
			levelStr = fl.Value.String()
		}

		fl = cmd.Flag("log-format")
		if fl != nil && (fl.Changed || logFormat == "") {
			logFormat = fl.Value.String()
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

	switch logFormat {
	case "plain":
		output := zerolog.ConsoleWriter{Out: out}
		output.FormatTimestamp = func(i interface{}) string { return "" }
		log.Logger = zerolog.New(output).Level(level).With().Logger()
	case "json":
		log.Logger = zerolog.New(out).Level(level)
	default:
		return fmt.Errorf("unknown log-format: %s", logFormat)
	}

	return nil
}
