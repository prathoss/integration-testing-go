package deps

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/prathoss/integration_testing/migrations"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func applyMigrations(t *testing.T, n *testcontainers.DockerNetwork, connStr string) {
	entries, err := migrations.FS.ReadDir(".")
	if err != nil {
		t.Fatalf("reading migrations fs: %v", err)
	}

	containerFiles := make([]testcontainers.ContainerFile, 0, len(entries))
	for _, entry := range entries {
		b, err := migrations.FS.ReadFile(entry.Name())
		if err != nil {
			t.Fatalf("reading migration %s: %v", entry.Name(), err)
		}
		containerFiles = append(
			containerFiles, testcontainers.ContainerFile{
				Reader:            bytes.NewReader(b),
				ContainerFilePath: filepath.Join("/", "migrations", entry.Name()),
			},
		)
	}
	migrateContainerReq := testcontainers.ContainerRequest{
		Image:      "migrate/migrate",
		Cmd:        []string{"-path=/migrations", "-database", connStr, "up"},
		WaitingFor: wait.ForExit(),
		Files:      containerFiles,
		Networks:   []string{n.Name},
	}
	container, err := testcontainers.GenericContainer(
		t.Context(), testcontainers.GenericContainerRequest{
			ContainerRequest: migrateContainerReq,
			Started:          true,
		},
	)
	if err != nil {
		t.Fatalf("starting migrations container: %v", err)
	}

	if err := container.Terminate(t.Context()); err != nil {
		t.Logf("failed to stop migrations container: %v", err)
	}
}
