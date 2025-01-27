package batron

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
)

var Version = "dev"
var Revision = "HEAD"

type GlobalOptions struct {
}

type CLI struct {
	Render  *RenderOption `cmd:"render" help:"Render job definition"`
	Deploy  *DeployOption `cmd:"deploy" help:"Deploy job definition"`
	Version VersionFlag   `name:"version" help:"show version"`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Printf("%s-%s\n", Version, Revision)
	app.Exit(0)
	return nil
}

func RunCLI(ctx context.Context, args []string) error {
	cli := CLI{
		Version: VersionFlag("0.1.0"),
	}
	parser, err := kong.New(&cli)
	if err != nil {
		return fmt.Errorf("error creating CLI parser: %w", err)
	}
	_, err = parser.Parse(args)
	if err != nil {
		fmt.Printf("error parsing CLI: %v\n", err)
		return fmt.Errorf("error parsing CLI: %w", err)
	}
	app, err := NewApp(&cli)
	if err != nil {
		return err
	}
	return app.Dispatch(ctx, args[0])
}
