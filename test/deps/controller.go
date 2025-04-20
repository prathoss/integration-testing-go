package deps

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type Controller struct {
	pg             *pgController
	profileService *serviceController
	pictureService *serviceController
}

func (c *Controller) Refresh(t *testing.T) {
	c.profileService.Stop(t)
	c.pictureService.Stop(t)

	c.pg.Restore(t)

	c.pictureService.Start(t)
	c.profileService.Start(t)
}

func (c *Controller) GetPG() *pgController {
	if c == nil {
		return nil
	}
	return c.pg
}

func (c *Controller) GetProfileService() *serviceController {
	if c == nil {
		return nil
	}
	return c.profileService
}

func (c *Controller) GetPictureService() *serviceController {
	if c == nil {
		return nil
	}
	return c.pictureService
}

type pgController struct {
	alias       string
	pgContainer *postgres.PostgresContainer
}

func (ctrl *pgController) Restore(t *testing.T) {
	if ctrl == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ctrl.pgContainer.Restore(ctx); err != nil {
		t.Fatalf("failed to restore postgres container: %v", err)
	}
}

func (ctrl *pgController) GetAddress(t *testing.T) string {
	if ctrl == nil {
		return ""
	}
	connStr, err := ctrl.pgContainer.ConnectionString(t.Context(), "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get postgres connection string: %v", err)
	}
	return connStr
}

func (ctrl *pgController) GetInternalAddress(t *testing.T) string {
	if ctrl == nil {
		return ""
	}
	publicAddr := ctrl.GetAddress(t)
	u, err := url.Parse(publicAddr)
	if err != nil {
		t.Fatalf("parsing connection string: %v", err)
	}
	u.Host = ctrl.alias
	return u.String()
}

type serviceController struct {
	alias     string
	port      nat.Port
	container testcontainers.Container
}

func (u *serviceController) Stop(t *testing.T) {
	if u == nil {
		return
	}
	if err := u.container.Stop(t.Context(), nil); err != nil {
		t.Fatalf("failed to stop service: %v", err)
	}
}

func (u *serviceController) Start(t *testing.T) {
	if u == nil {
		return
	}
	if err := u.container.Start(t.Context()); err != nil {
		t.Fatalf("failed to start service: %v", err)
	}
}

func (u *serviceController) GetAddress(t *testing.T) string {
	if u == nil {
		return ""
	}

	port, err := u.container.MappedPort(t.Context(), u.port)
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}

	host, err := u.container.Host(t.Context())
	if err != nil {
		t.Fatalf("failed to get host: %v", err)
	}

	return fmt.Sprintf("http://%s:%s", host, port.Port())
}

func (u *serviceController) GetInternalAddress() string {
	if u == nil {
		return ""
	}
	return fmt.Sprintf("http://%s:%s", u.alias, u.port.Port())
}
