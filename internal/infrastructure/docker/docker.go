package docker

import (
	"context"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Docker struct {
	cli *client.Client
}

func NewDocker() (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Docker{
		cli: cli,
	}, nil
}

func (d *Docker) RunContainer(image string, args []string, envVars map[string]string) (string, error) {
	config := &container.Config{
		Image: image,
		Cmd:   args,
		Env:   convertEnvVarsToStringSlice(envVars),
	}

	resp, err := d.cli.ContainerCreate(context.Background(), config, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := d.cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (d *Docker) GetContainerLogs(containerID string) (string, error) {
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}

	out, err := d.cli.ContainerLogs(context.Background(), containerID, options)
	if err != nil {
		return "", err
	}
	defer out.Close()

	var logsBuilder strings.Builder
	_, err = io.Copy(&logsBuilder, out)
	if err != nil {
		return "", err
	}

	return logsBuilder.String(), nil
}

func convertEnvVarsToStringSlice(envVars map[string]string) []string {
	var env []string
	for key, value := range envVars {
		env = append(env, key+"="+value)
	}
	return env
}
