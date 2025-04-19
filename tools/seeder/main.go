package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
	"github.com/prathoss/integration_testing/seeder"
)

func main() {
	path := flag.String("path", "", "path that contains database seeds")
	uri := flag.String("uri", "", "uri to connect to database")

	flag.Parse()

	var err error = nil
	if *path == "" {
		err = errors.Join(errors.New("path is required"))
	}
	if *uri == "" {
		err = errors.Join(errors.New("uri is required"))
	}

	if err != nil {
		fmt.Printf("%s\n\n", err.Error())
		flag.Usage()
		os.Exit(1)
	}

	err = func() error {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		conn, err := pgx.Connect(ctx, *uri)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return err
		}

		dirFS := os.DirFS(*path)
		rdfs, ok := dirFS.(fs.ReadDirFS)
		if !ok {
			fmt.Printf("%s\n", "path is not a directory")
			return err
		}

		if err := seeder.Seed(ctx, conn, rdfs); err != nil {
			fmt.Printf("%s\n", err.Error())
			return err
		}

		return nil
	}()
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("Seeded successfully")
}
