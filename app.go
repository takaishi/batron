package batron

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/batch"
)

type App struct {
	CLI         *CLI
	batchClient *batch.Client
}

func NewApp(cli *CLI) (*App, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config, %w", err)
	}
	return &App{
		CLI:         cli,
		batchClient: batch.NewFromConfig(awsConfig),
	}, nil
}

func (a *App) Dispatch(ctx context.Context, command string) error {
	switch command {
	case "render":
		cmd, err := NewRenderCommand(a, a.CLI.Render)
		if err != nil {
			return err
		}
		return cmd.Run(ctx)
	case "deploy":
		cmd, err := NewDeployCommand(a, a.CLI.Deploy)
		if err != nil {
			return err
		}
		return cmd.Run(ctx)
	case "deregister":
		cmd, err := NewDeregisterCommand(a, a.CLI.Deregister)
		if err != nil {
			return err
		}
		return cmd.Run(ctx)
	}
	return nil
}
