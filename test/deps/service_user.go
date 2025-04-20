package deps

import (
	"context"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func (b *builder) WithProfileService(t *testing.T) *builder {
	if b.ctrl.profileService != nil {
		return b
	}
	if b.ctrl.pg == nil {
		b.WithPG(t)
	}

	alias := "profile"
	port := nat.Port("8080/tcp")
	profileServerRequest := testcontainers.ContainerRequest{
		Image: "gopic/profile:under-test",
		Env: map[string]string{
			"GOPIC_DATABASE_URI": b.ctrl.pg.GetInternalAddress(t),
			"GOPIC_SERVER_PORT":  "8080",
		},
		ExposedPorts: []string{string(port)},
		Networks:     []string{b.network.Name},
		NetworkAliases: map[string][]string{
			b.network.Name: {alias},
		},
		WaitingFor: wait.ForExposedPort(),
	}

	container, err := testcontainers.GenericContainer(
		t.Context(), testcontainers.GenericContainerRequest{
			ContainerRequest: profileServerRequest,
			Started:          true,
		},
	)
	if err != nil {
		t.Fatalf("starting profile service container: %v", err)
	}

	t.Cleanup(
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := container.Terminate(ctx); err != nil {
				t.Logf("failed to stop profile service container: %v", err)
			}
		},
	)

	ctrl := &serviceController{
		alias:     alias,
		port:      port,
		container: container,
	}
	b.ctrl.profileService = ctrl
	return b
}
