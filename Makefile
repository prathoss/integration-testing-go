fmt:
	golangci-lint fmt

lint: fmt
	golangci-lint run

build-for-testing-picture:
	docker build -f service_picture/Dockerfile . -t gopic/picture:under-test

build-for-testing-profile:
	docker build -f service_profile/Dockerfile . -t gopic/profile:under-test

build-for-testing: build-for-testing-picture build-for-testing-profile

test: build-for-test
	go test -v ./...

create-migration: check-migration-name
	migrate create -ext sql -dir migrations -seq ${MIGRATION_NAME}

migrate: check-conn
	migrate -database "${DATABASE_URL}" -path migrations up

seed: check-conn
	go run ./tools/seeder/ -path ./seed -uri "${DATABASE_URL}"

check-conn:
ifndef DATABASE_URL
	$(error DATABASE_URL is not set)
endif

check-migration-name:
ifndef MIGRATION_NAME
	$(error MIGRATION_NAME is not set)
endif
