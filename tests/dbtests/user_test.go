package dbtests

import (
	"context"
	"tefsi/internal/domain"
	"tefsi/tests"
	"testing"
)

func TestCreateUser(t *testing.T) {
	container, db, err := tests.CreateContainer("test-db")
	defer db.Close()
	defer container.Terminate(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	repos, err := tests.CreateRepos(db)
	if err != nil {
		t.Fatal(err)
	}

	err = repos.UserRepository.CreateUser(context.Background(), &domain.User{Login: "user1", Password: "password", IsAdmin: false})
	if err != nil {
		t.Fatal(err)
	}
}
