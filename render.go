package batron

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/google/go-jsonnet"
)

type RenderCommand struct {
	batchClient  *batch.Client
	RenderOption *RenderOption
}

type RenderOption struct {
	Config string            `help:"Config file"`
	ExtStr map[string]string `help:"ExtVar" set:"Key=Val"`
}

func NewRenderCommand(app *App, option *RenderOption) (*RenderCommand, error) {
	return &RenderCommand{
		batchClient:  app.batchClient,
		RenderOption: option,
	}, nil
}

func (c *RenderCommand) Run(ctx context.Context) error {
	config, err := c.loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	jobDef, err := c.renderJobDefinition(config)
	if err != nil {
		return fmt.Errorf("failed to render job definition: %w", err)
	}
	json, err := json.MarshalIndent(jobDef, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal job definition: %w", err)
	}
	fmt.Printf("%s\n", string(json))
	return nil
}

func (c *RenderCommand) configFileDir() string {
	return filepath.Dir(c.RenderOption.Config)
}

func (c *RenderCommand) loadConfig() (*Config, error) {
	vm := jsonnet.MakeVM()
	for k, v := range c.RenderOption.ExtStr {
		vm.ExtVar(k, v)
	}
	data, err := vm.EvaluateFile(c.RenderOption.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate jsonnet: %w", err)
	}

	var config Config
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &config, nil
}

func (c *RenderCommand) renderJobDefinition(config *Config) (*batch.RegisterJobDefinitionInput, error) {
	vm := jsonnet.MakeVM()
	for k, v := range c.RenderOption.ExtStr {
		vm.ExtVar(k, v)
	}
	data, err := vm.EvaluateFile(filepath.Join(c.configFileDir(), config.JobDefinition))
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate jsonnet: %w", err)
	}
	var jobDef batch.RegisterJobDefinitionInput
	if err := json.Unmarshal([]byte(data), &jobDef); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job definition: %w", err)
	}
	return &jobDef, nil
}
