package batron

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/aws/aws-sdk-go-v2/service/batch/types"
)

type SubmitJobCommand struct {
	*App
	SubmitJobOption *SubmitJobOption
}

type SubmitJobOption struct {
	JobDefinition string            `help:"Job definition name" required:""`
	JobQueue      string            `help:"Job queue name" required:""`
	JobName       string            `help:"Job name" required:""`
	Command       []string          `help:"Command to run"`
	Environment   map[string]string `help:"Environment variables" set:"Key=Val"`
}

func NewSubmitJobCommand(app *App, option *SubmitJobOption) (*SubmitJobCommand, error) {
	return &SubmitJobCommand{
		App:             app,
		SubmitJobOption: option,
	}, nil
}

func (c *SubmitJobCommand) Run(ctx context.Context) error {
	var envVars []types.KeyValuePair
	for k, v := range c.SubmitJobOption.Environment {
		name := k
		value := v
		envVars = append(envVars, types.KeyValuePair{
			Name:  &name,
			Value: &value,
		})
	}

	input := &batch.SubmitJobInput{
		JobDefinition: aws.String(c.SubmitJobOption.JobDefinition),
		JobName:       aws.String(c.SubmitJobOption.JobName),
		JobQueue:      aws.String(c.SubmitJobOption.JobQueue),
	}

	if len(c.SubmitJobOption.Command) > 0 {
		input.ContainerOverrides = &types.ContainerOverrides{
			Command:     c.SubmitJobOption.Command,
			Environment: envVars,
		}
	}

	output, err := c.batchClient.SubmitJob(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to submit job: %w", err)
	}

	outputJson, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}
	fmt.Println(string(outputJson))

	return nil
}
