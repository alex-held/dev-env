package sdk

import (
	"io"

	"github.com/gobuffalo/plugins/plugflag"
	"github.com/spf13/pflag"

	flagger2 "github.com/alex-held/devctl/pkg/plugins/flagger"
)

func (cmd *Cmd) Flags() *pflag.FlagSet {
	if cmd.flags != nil {
		return cmd.flags
	}
	flags := pflag.NewFlagSet(cmd.PluginName(), pflag.ContinueOnError)
	flags.BoolVarP(&cmd.help, "help", "h", false, "print this help")

	for _, p := range cmd.ScopedPlugins() {
		switch t := p.(type) {
		case Flagger:
			for _, f := range plugflag.Clean(p, t.SetupFlags()) {
				flags.AddGoFlag(f)
			}
		case Pflagger:
			for _, f := range flagger2.CleanPflags(p, t.SetupFlags()) {
				flags.AddFlag(f)
			}
		}
	}

	cmd.flags = flags
	return cmd.flags
}

func (cmd *Cmd) PrintFlags(w io.Writer) error {
	flags := cmd.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}
