package seeder

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5"
)

func Seed(ctx context.Context, conn *pgx.Conn, f fs.ReadDirFS) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	err = func(ctx context.Context, conn pgx.Tx) error {
		entries, err := f.ReadDir(".")
		if err != nil {
			return fmt.Errorf("reading directory %s: %w", ".", err)
		}

		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".csv" {
				continue
			}

			fileName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			_, tableName, found := strings.Cut(fileName, "_")
			if !found {
				tableName = fileName
			}
			file, err := f.Open(entry.Name())
			if err != nil {
				return fmt.Errorf("reading file %s: %w", entry.Name(), err)
			}

			if _, err := conn.Conn().PgConn().CopyFrom(
				ctx,
				file,
				fmt.Sprintf("COPY %s FROM STDIN (FORMAT csv, DELIMITER ',', HEADER)", tableName),
			); err != nil {
				return fmt.Errorf("copying file %s: %w", entry.Name(), err)
			}

			if err := file.Close(); err != nil {
				return fmt.Errorf("closing file %s: %w", entry.Name(), err)
			}
		}
		return nil
	}(ctx, tx)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			err = errors.Join(err, fmt.Errorf("rollback transaction: %w", txErr))
		}
		return err
	}
	return tx.Commit(ctx)
}
