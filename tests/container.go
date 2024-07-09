package tests

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreateContainer(dbName string) (testcontainers.Container, *pgxpool.Pool, error) {
	port := "5432/tcp"

	env := map[string]string{
		"POSTGRES_USER":     dbName,
		"POSTGRES_DB":       dbName,
		"POSTGRES_PASSWORD": "password",
	}

	req := testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres",
			ExposedPorts: []string{port},
			Cmd:          []string{"postgres"},
			Env:          env,
			Name:         dbName + uuid.New().String(),
			WaitingFor: wait.ForSQL(nat.Port(port), "pgx", func(host string, port nat.Port) string {
				return fmt.Sprintf("postgres://%s:password@%s:%s/%s?sslmode=disable", dbName, host, port.Port(), dbName)
			}),
		},
	}

	container, err := testcontainers.GenericContainer(context.Background(), req)
	if err != nil {
		return nil, nil, err
	}

	mappedPort, err := container.MappedPort(context.Background(), nat.Port(port))
	if err != nil {
		return nil, nil, err
	}

	url := fmt.Sprintf("postgres://%s:password@%s:%s/%s?sslmode=disable", dbName, "localhost", mappedPort.Port(), dbName)
	db, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, nil, err
	}

	return container, db, nil
}
