package deps

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/prathoss/integration_testing/seed"
	"github.com/prathoss/integration_testing/seeder"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

func (b *builder) WithPG(t *testing.T) *builder {
	if b.ctrl.pg != nil {
		return b
	}

	alias := "postgres"
	postgresContainer, err := postgres.Run(
		t.Context(),
		"postgres:17-alpine",
		postgres.WithDatabase("gopic"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithSQLDriver("pgx"),
		network.WithNetwork([]string{alias}, b.network),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
		),
	)
	if err != nil {
		t.Fatalf("starting postgres test container: %v", err)
	}
	t.Cleanup(
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := postgresContainer.Terminate(ctx); err != nil {
				t.Logf("failed to terminate postgres container: %v", err)
			}
		},
	)

	ctrl := &pgController{
		alias:       alias,
		pgContainer: postgresContainer,
	}

	// apply migrations
	applyMigrations(t, b.network, ctrl.GetInternalAddress(t))

	conn, err := pgx.Connect(t.Context(), ctrl.GetAddress(t))
	if err != nil {
		t.Fatalf("connecting to postgres: %v", err)
	}

	// seed database
	if err := seeder.Seed(t.Context(), conn, seed.FS); err != nil {
		t.Fatalf("seeding db: %v", err)
	}

	// reset sequences after COPY
	if _, err := conn.Exec(
		t.Context(),
		"SELECT setval('profiles_id_seq', max(id)) FROM profiles;",
	); err != nil {
		t.Fatalf("could not reset sequence for profiles: %v", err)
	}
	if _, err := conn.Exec(
		t.Context(),
		"SELECT setval('pictures_id_seq', max(id)) FROM pictures;",
	); err != nil {
		t.Fatalf("could not reset sequence for pictures: %v", err)
	}
	if _, err := conn.Exec(
		t.Context(),
		"SELECT setval('picture_views_id_seq', max(id)) FROM picture_views;",
	); err != nil {
		t.Fatalf("could not reset sequence for picture_views: %v", err)
	}

	// close connection to be able to create snapshot
	if err := conn.Close(t.Context()); err != nil {
		t.Fatalf("closing connection: %v", err)
	}

	// create snapshot of the container eg savepoint to return to
	if err := postgresContainer.Snapshot(t.Context()); err != nil {
		t.Fatalf("snapshotting postgres container: %v", err)
	}

	b.ctrl.pg = ctrl
	return b
}
