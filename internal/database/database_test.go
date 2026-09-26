package database

import (
	"database/sql"
	"os"
	"testing"
)

func testDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("KDOCS_DB_DSN")

	if dsn == "" {
		dsn = "postgres://admin:Incorrect@localhost:5432/kdocs?sslmode=disable"
	}
	db, err := Open(dsn)

	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	return db
}

func TestFindSpace(t *testing.T) {
	db := testDB(t)

	space, found, err := FindSpace(db, "databases")
	if err != nil {
		t.Fatalf("FindSpace returned an error: %v", err)
	}

	if !found {
		t.Fatal("expected databases space to be found")
	}

	if space.Name != "Databases" {
		t.Errorf("expected Databases, got %q", space.Name)
	}
}

func TestFindSpaceNotFound(t *testing.T) {
	db := testDB(t)

	_, found, err := FindSpace(db, "does-not-exist")
	if err != nil {
		t.Fatalf("FindSpace returned an error: %v", err)
	}

	if found {
		t.Fatal("expected space not to be found")
	}
}

func TestFindPagesBySpace(t *testing.T) {
	db := testDB(t)

	pages, err := FindPagesBySpace(db, "databases")
	if err != nil {
		t.Fatalf("FindPagesBySpace returned an error: %v", err)
	}

	if len(pages) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(pages))
	}
}

func TestFindPage(t *testing.T) {
	db := testDB(t)

	page, found, err := FindPage(db, "security-baseline")
	if err != nil {
		t.Fatalf("FindPage returned an error: %v", err)
	}

	if !found {
		t.Fatal("expected security-baseline page to be found")
	}

	if page.Title != "Application Security Baseline" {
		t.Errorf(
			"expected Application Security Baseline, got %q",
			page.Title,
		)
	}

	if page.SpaceSlug != "security" {
		t.Errorf(
			"expected security space, got %q",
			page.SpaceSlug,
		)
	}
}
