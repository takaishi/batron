package batron

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/google/go-jsonnet"
)

type DeployCommand struct {
	*App
	DeployOption *DeployOption
}

type DeployOption struct {
	Config string            `help:"Config file"`
	ExtStr map[string]string `help:"ExtVar" set:"Key=Val"`
}

func NewDeployCommand(app *App, option *DeployOption) (*DeployCommand, error) {
	return &DeployCommand{
		App:          app,
		DeployOption: option,
	}, nil
}

func (c *DeployCommand) Run(ctx context.Context) error {
	config, err := c.loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	jobDef, err := c.renderJobDefinition(config)
	if err != nil {
		return fmt.Errorf("failed to render job definition: %w", err)
	}

	output, err := c.batchClient.RegisterJobDefinition(ctx, jobDef)
	if err != nil {
		return fmt.Errorf("failed to register job definition: %w", err)
	}
	outputJson, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}
	fmt.Println(string(outputJson))

	return nil
}
func (c *DeployCommand) loadConfig() (*Config, error) {
	vm := jsonnet.MakeVM()
	for k, v := range c.DeployOption.ExtStr {
		vm.ExtVar(k, v)
	}
	data, err := vm.EvaluateFile(c.DeployOption.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate jsonnet: %w", err)
	}

	var config Config
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &config, nil
}
func (c *DeployCommand) configFileDir() string {
	return filepath.Dir(c.DeployOption.Config)
}

func (c *DeployCommand) renderJobDefinition(config *Config) (*batch.RegisterJobDefinitionInput, error) {
	vm := jsonnet.MakeVM()
	for k, v := range c.DeployOption.ExtStr {
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
