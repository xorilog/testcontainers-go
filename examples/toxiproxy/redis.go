package toxiproxy

import (
	"context"

	dockerc "github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type redisContainer struct {
	testcontainers.Container
}

func setupRedis(ctx context.Context, network string, networkAlias []string) (*redisContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:6",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
		Networks: []string{
			network,
		},
		NetworkAliases: map[string][]string{
			network: networkAlias,
		},

		HostConfigModifier: func(config *dockerc.HostConfig) {
			config.NetworkMode = "slirp4netns"
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &redisContainer{Container: container}, nil
}
