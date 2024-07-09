package tests

import (
	"context"
	"testing"
)

func TestCreateContainer(t *testing.T) {
	container, db, err := CreateContainer("test-db")
	if err != nil {
		t.Fatal("couldnt create container:\n", err)
	}
	defer db.Close()
	defer container.Terminate(context.Background())

	_, err = db.Exec(context.Background(), "CREATE TABLE testtable ()")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateRepos(t *testing.T) {
	container, db, err := CreateContainer("test-db")
	if err != nil {
		t.Fatal("couldnt create container:\n", err)
	}
	defer db.Close()
	defer container.Terminate(context.Background())

	_, err = CreateRepos(db)
	if err != nil {
		t.Fatal(err)
	}
}
