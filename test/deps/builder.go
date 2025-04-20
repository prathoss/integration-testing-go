package deps

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/network"
)

type builder struct {
	network *testcontainers.DockerNetwork
	ctrl    *Controller
}

func NewBuilder(t *testing.T) *builder {
	disable := os.Getenv("TESTCONTAINERS_DISABLE_LOGGING")
	if disable == "true" {
		testcontainers.DefaultLoggingHook = func(logger log.Logger) testcontainers.ContainerLifecycleHooks {
			return testcontainers.ContainerLifecycleHooks{}
		}
	}

	n, err := network.New(t.Context())
	if err != nil {
		t.Fatalf("failed to create network: %v", err)
	}

	t.Cleanup(
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := n.Remove(ctx)
			if err != nil {
				t.Logf("failed to remove network: %v", err)
			}
		},
	)

	return &builder{
		network: n,
		ctrl:    &Controller{},
	}
}

func (b *builder) Build() *Controller {
	return b.ctrl
}
