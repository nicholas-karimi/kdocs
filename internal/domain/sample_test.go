package domain

import "testing"

func TestFindSpace(t *testing.T) {
	space, found := FindSpace("databases")
	if !found {
		t.Fatal("expected databases space to be found")
	}

	if space.Name != "Databases" {
		t.Errorf("expected Databases, got %q", space.Name)
	}
}

func TestFindSpaceNotFound(t *testing.T) {
	_, found := FindSpace("unknown")

	if found {
		t.Fatal("expected space not to be found")
	}
}

func TestFindPage(t *testing.T) {
	page, found := FindPage("postgresql-migration")

	if !found {
		t.Fatal("expected page to be found")
	}

	if page.Title != "PostgreSQL Migration" {
		t.Errorf("expected PostgreSQL Migration, got %q", page.Title)
	}
}

func TestFindPageNotFound(t *testing.T) {
	_, found := FindPage("unknown")

	if found {
		t.Fatal("expected page not to be found")
	}
}

func TestFindPagesBySpace(t *testing.T) {
	pages := FindPagesBySpace("databases")

	if len(pages) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(pages))
	}
}

func TestFindPagesByUnknownSpace(t *testing.T) {
	pages := FindPagesBySpace("unknown")

	if len(pages) != 0 {
		t.Fatalf("expected 0 pages, got %d", len(pages))
	}
}
