package plugin

import (
	"github.com/maliceio/cli/opts"
	"github.com/maliceio/engine/cli/malice"
	"github.com/maliceio/engine/cli/malice/command"
	"github.com/maliceio/engine/cli/malice/command/formatter"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type listOptions struct {
	quiet   bool
	noTrunc bool
	format  string
	filter  opts.FilterOpt
}

func newListCommand(maliceCli *command.MaliceCli) *cobra.Command {
	opts := listOptions{filter: opts.NewFilterOpt()}

	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Short:   "List plugins",
		Aliases: []string{"list"},
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(maliceCli, opts)
		},
	}

	flags := cmd.Flags()

	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only display plugin IDs")
	flags.BoolVar(&opts.noTrunc, "no-trunc", false, "Don't truncate output")
	flags.StringVar(&opts.format, "format", "", "Pretty-print plugins using a Go template")
	flags.VarP(&opts.filter, "filter", "f", "Provide filter values (e.g. 'enabled=true')")

	return cmd
}

func runList(maliceCli *command.MaliceCli, opts listOptions) error {
	plugins, err := maliceCli.Client().PluginList(context.Background(), opts.filter.Value())
	if err != nil {
		return err
	}

	format := opts.format
	if len(format) == 0 {
		if len(maliceCli.ConfigFile().PluginsFormat) > 0 && !opts.quiet {
			format = maliceCli.ConfigFile().PluginsFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	pluginsCtx := formatter.Context{
		// Output: maliceCli.Out(),
		Format: formatter.NewPluginFormat(format, opts.quiet),
		Trunc:  !opts.noTrunc,
	}
	return formatter.PluginWrite(pluginsCtx, plugins)
}
