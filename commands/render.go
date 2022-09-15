package commands

import (
	"encoding"
	"encoding/json"
	"io"

	"gopkg.in/yaml.v3"
)

func render(w io.Writer, format string, v interface{}) error {
	if format == "yaml" || format == "yml" {
		return renderYAML(w, v)
	}

	if format == "json" {
		return renderJSON(w, v)
	}

	if wt, ok := v.(io.WriterTo); ok {
		_, err := wt.WriteTo(w)

		return err
	}

	if tm, ok := v.(encoding.TextMarshaler); ok {
		b, err := tm.MarshalText()
		if err != nil {
			return err
		}

		_, err = w.Write(b)

		return err
	}

	return renderJSON(w, v)
}

func renderYAML(w io.Writer, v interface{}) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)

	return enc.Encode(v)
}

func renderJSON(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	return enc.Encode(v)
}
