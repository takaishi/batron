package batron

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/google/go-jsonnet"
)

type DeregisterCommand struct {
	*App
	DeregisterOption *DeregisterOption
}

type DeregisterOption struct {
	Config string            `help:"Config file"`
	ExtStr map[string]string `help:"ExtStr" set:"Key=Val"`
	Keeps  int               `help:"Number of job definition generations to keep" default:"5"`
}

func NewDeregisterCommand(app *App, option *DeregisterOption) (*DeregisterCommand, error) {
	return &DeregisterCommand{
		App:              app,
		DeregisterOption: option,
	}, nil
}

func (c *DeregisterCommand) Run(ctx context.Context) error {
	config, err := c.loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	jobDef, err := c.renderJobDefinition(config)
	if err != nil {
		return fmt.Errorf("failed to render job definition: %w", err)
	}

	if err := c.cleanupOldJobDefinitions(ctx, *jobDef.JobDefinitionName, c.DeregisterOption.Keeps); err != nil {
		return fmt.Errorf("failed to cleanup old job definitions: %w", err)
	}

	return nil
}

func (c *DeregisterCommand) loadConfig() (*Config, error) {
	vm := jsonnet.MakeVM()
	for k, v := range c.DeregisterOption.ExtStr {
		vm.ExtVar(k, v)
	}
	data, err := vm.EvaluateFile(c.DeregisterOption.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate jsonnet: %w", err)
	}

	var config Config
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &config, nil
}

func (c *DeregisterCommand) configFileDir() string {
	return filepath.Dir(c.DeregisterOption.Config)
}

func (c *DeregisterCommand) renderJobDefinition(config *Config) (*batch.RegisterJobDefinitionInput, error) {
	vm := jsonnet.MakeVM()
	for k, v := range c.DeregisterOption.ExtStr {
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

func (c *DeregisterCommand) cleanupOldJobDefinitions(ctx context.Context, jobDefName string, keeps int) error {
	input := &batch.DescribeJobDefinitionsInput{
		JobDefinitionName: &jobDefName,
		Status:            aws.String("ACTIVE"),
	}

	output, err := c.batchClient.DescribeJobDefinitions(ctx, input)
	if err != nil {
		return err
	}

	sort.Slice(output.JobDefinitions, func(i, j int) bool {
		return *output.JobDefinitions[i].Revision > *output.JobDefinitions[j].Revision
	})

	for i := keeps; i < len(output.JobDefinitions); i++ {
		deregisterInput := &batch.DeregisterJobDefinitionInput{
			JobDefinition: output.JobDefinitions[i].JobDefinitionArn,
		}

		if _, err := c.batchClient.DeregisterJobDefinition(ctx, deregisterInput); err != nil {
			return err
		}
		fmt.Printf("Deregistered job definition: %s\n", *output.JobDefinitions[i].JobDefinitionArn)
	}

	return nil
}
