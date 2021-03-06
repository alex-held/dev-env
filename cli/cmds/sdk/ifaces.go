package sdk

import (
	"context"
	"flag"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/spf13/pflag"
)

type Installer interface {
	plugins.Plugin
	Install(sdk string) error
}

type Sdker interface {
	plugins.Plugin
	ExecuteCommand(ctx context.Context, root string, args []string) error
}

type Namer interface {
	plugins.Plugin
	CmdName() string
}

type Flagger interface {
	plugins.Plugin
	SetupFlags() []*flag.Flag
}

type Pflagger interface {
	plugins.Plugin
	SetupFlags() []*pflag.Flag
}

type Aliaser interface {
	Sdker
	CmdAliases() []string
}

type Stdouter = plugio.Outer
type Needer = plugins.Needer
