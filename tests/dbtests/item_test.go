package dbtests

import (
	"context"
	"log"
	"tefsi/internal/domain"
	"tefsi/tests"
	"testing"
)

// tests item equality without IDs
func itemEq(item1 domain.Item, item2 domain.Item) bool {
	item1.ID = 0
	item2.ID = 0
	return item1 == item2
}

// TODO: GetItemsByID, DeleteItem

func TestCreateItem(t *testing.T) {
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

	allCats, err := repos.CategoryRepository.GetCategories(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	categoryID, err := categoryIDFromTitle("cat", *allCats)
	if err != nil {
		t.Fatal(err)
	}

	repos.ItemRepository.CreateItem(context.Background(), &domain.Item{
		Title:         "item1",
		Description:   "very cool item !",
		Price:         999,
		CategoryID:    categoryID,
		CategoryTitle: "cat",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetItems(t *testing.T) {
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

	allCats, err := repos.CategoryRepository.GetCategories(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	catID, err := categoryIDFromTitle("cat", *allCats)
	if err != nil {
		t.Fatal(err)
	}
	carID, err := categoryIDFromTitle("car", *allCats)
	if err != nil {
		t.Fatal(err)
	}

	cat1 := domain.Item{
		Title:         "cat1",
		Description:   "cat is cat because CATegory",
		Price:         999,
		CategoryID:    catID,
		CategoryTitle: "cat",
	}
	cat2 := domain.Item{
		Title:         "cat2",
		Description:   "meow",
		Price:         12,
		CategoryID:    catID,
		CategoryTitle: "cat",
	}
	car1 := domain.Item{
		Title:         "mashina1",
		Description:   "car         !",
		Price:         123123,
		CategoryID:    carID,
		CategoryTitle: "car",
	}

	repos.ItemRepository.CreateItem(context.Background(), &cat1)
	if err != nil {
		t.Fatal(err)
	}
	repos.ItemRepository.CreateItem(context.Background(), &cat2)
	if err != nil {
		t.Fatal(err)
	}
	repos.ItemRepository.CreateItem(context.Background(), &car1)
	if err != nil {
		t.Fatal(err)
	}

	filterCat := domain.Filter{
		CategoryID: catID,
	}
	filterCar := domain.Filter{
		CategoryID: carID,
	}
	filter1 := domain.Filter{
		SearchString: "1",
	}
	filterMashina := domain.Filter{
		SearchString: "mashina",
	}
	filterAll := domain.Filter{}

	catItems, err := repos.ItemRepository.GetItems(context.Background(), &filterCat)
	if err != nil {
		t.Fatal(err)
	}
	carItems, err := repos.ItemRepository.GetItems(context.Background(), &filterCar)
	if err != nil {
		t.Fatal(err)
	}
	oneItems, err := repos.ItemRepository.GetItems(context.Background(), &filter1)
	if err != nil {
		log.Println(filter1.GenerateString())
		t.Fatal(err)
	}
	mashinaItems, err := repos.ItemRepository.GetItems(context.Background(), &filterMashina)
	if err != nil {
		t.Fatal(err)
	}
	allItems, err := repos.ItemRepository.GetItems(context.Background(), &filterAll)

	if len(*catItems) != 2 {
		t.Fatal("expected 2 cat items, got", len(*catItems))
	}
	if !itemEq((*catItems)[0], cat1) || !itemEq((*catItems)[1], cat2) {
		t.Fatalf("expected %v and %v, got %v and %v", cat1, cat2, (*catItems)[0], (*catItems)[1])
	}

	if len(*carItems) != 1 {
		t.Fatal("expected 1 car item, got", len(*carItems))
	}
	if !itemEq((*carItems)[0], car1) {
		t.Fatalf("expected %v, got %v", car1, (*carItems)[0])
	}

	if len(*oneItems) != 2 {
		t.Fatal("expected 2 oneitems, got", len(*oneItems))
	}
	if !itemEq((*oneItems)[0], cat1) || !itemEq((*oneItems)[1], car1) {
		t.Fatalf("expected %v and %v, got %v and %v", cat1, car1, (*oneItems)[0], (*oneItems)[1])
	}

	if len(*mashinaItems) != 1 {
		t.Fatal("expected 1 mashina item, got", len(*mashinaItems))
	}
	if !itemEq((*mashinaItems)[0], car1) {
		t.Fatalf("expected %v, got %v", car1, (*mashinaItems)[0])
	}

	if len(*allItems) != 3 {
		t.Fatal("expected 3 items, got", len(*allItems))
	}
	// if (*allItems)[0] != cat1 || (*allItems)[1] != cat2 || (*allItems)[2] != car1 {
	if !itemEq((*allItems)[0], cat1) || !itemEq((*allItems)[1], cat2) || !itemEq((*allItems)[2], car1) {
		t.Fatalf("expected %v, %v and %v, got %v, %v and %v", cat1, cat2, car1, (*allItems)[0], (*allItems)[1], (*allItems)[2])
	}
}
