package nginx

import (
	"context"
	"fmt"
	"time"

	dockerc "github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type nginxContainer struct {
	testcontainers.Container
	URI string
}

func startContainer(ctx context.Context) (*nginxContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "docker.io/nginxinc/nginx-unprivileged",
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForHTTP("/").WithStartupTimeout(10 * time.Second),
		Networks: []string{
			"slirp4netns",
		},
		HostConfigModifier: func(config *dockerc.HostConfig) {
			config.NetworkMode = "slirp4netns"
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ProviderType:     testcontainers.ProviderPodman,
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "8080")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &nginxContainer{Container: container, URI: uri}, nil
}
