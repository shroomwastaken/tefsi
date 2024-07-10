package dbtests

import (
	"context"
	"fmt"
	"tefsi/internal/domain"
	"tefsi/tests"
	"testing"
)

func categoryIDFromTitle(title string, categories []domain.Category) (int, error) {
	var id int
	found := false

	for _, category := range categories {
		if category.Title != title {
			continue
		}

		if found {
			return 0, fmt.Errorf("There are multiple categories with the same title which is like fine but id from title doesnt know what to do")
		}
		id = category.ID
		found = true
	}

	return id, nil
}

func TestCreateCategory(t *testing.T) {
	container, db, err := tests.CreateContainer("test-db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer container.Terminate(context.Background())

	repos, err := tests.CreateRepos(db)
	if err != nil {
		t.Fatal(err)
	}

	err = repos.CategoryRepository.CreateCategory(context.Background(), &domain.Category{
		Title: "cat",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetCategories(t *testing.T) {
	container, db, err := tests.CreateContainer("test-db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer container.Terminate(context.Background())

	repos, err := tests.CreateRepos(db)
	if err != nil {
		t.Fatal(err)
	}

	err = repos.CategoryRepository.CreateCategory(context.Background(), &domain.Category{
		Title: "cat",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = repos.CategoryRepository.CreateCategory(context.Background(), &domain.Category{
		Title: "car",
	})
	if err != nil {
		t.Fatal(err)
	}

	catsp, err := repos.CategoryRepository.GetCategories(context.Background())
	cats := *catsp
	if err != nil {
		t.Fatal(err)
	}

	if len(cats) != 2 {
		t.Fatal("expected 2 categories, got", len(cats))
	}

	if cats[0].Title != "cat" || cats[1].Title != "car" {
		t.Fatal("expected cat and car, got", cats[0].Title, "and", cats[1].Title)
	}
}

func TestGetCategoryByID(t *testing.T) {
	container, db, err := tests.CreateContainer("test-db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer container.Terminate(context.Background())

	repos, err := tests.CreateRepos(db)
	if err != nil {
		t.Fatal(err)
	}

	err = repos.CategoryRepository.CreateCategory(context.Background(), &domain.Category{
		Title: "cat",
	})
	if err != nil {
		t.Fatal(err)
	}

	catsp, err := repos.CategoryRepository.GetCategories(context.Background())
	cats := *catsp
	if err != nil {
		t.Fatal(err)
	}
	if len(cats) != 1 {
		t.Fatal("expected 1 category, got", len(cats))
	}

	id := cats[0].ID
	cat, err := repos.CategoryRepository.GetCategoryByID(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}

	if *cat != cats[0] {
		t.Fatalf("expected 2 equal categories, got: %+v and %+v", *cat, cats[0])
	}
}

func TestDeleteCategory(t *testing.T) {
	container, db, err := tests.CreateContainer("test-db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer container.Terminate(context.Background())

	repos, err := tests.CreateRepos(db)
	if err != nil {
		t.Fatal(err)
	}

	err = repos.CategoryRepository.CreateCategory(context.Background(), &domain.Category{
		Title: "cat",
	})
	if err != nil {
		t.Fatal(err)
	}

	catsp, err := repos.CategoryRepository.GetCategories(context.Background())
	cats := *catsp
	if err != nil {
		t.Fatal(err)
	}
	if len(cats) != 1 {
		t.Fatal("expected 1 category, got", len(cats))
	}

	id := cats[0].ID
	err = repos.CategoryRepository.DeleteCategory(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}

	catsAfterDel, err := repos.CategoryRepository.GetCategories(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(*catsAfterDel) != 0 {
		t.Fatal("the category wasnt deleted")
	}
}
