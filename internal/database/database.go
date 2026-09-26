package database

import (
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nicholas-karimi/kdocs/internal/domain"
)

func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// single space
func FindSpace(db *sql.DB, slug string) (domain.Space, bool, error) {
	var space domain.Space

	err := db.QueryRow(
		"SELECT name, slug FROM spaces WHERE slug = $1",
		slug,
	).Scan(&space.Name, &space.Slug)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Space{}, false, nil
		}
		return domain.Space{}, false, err

	}
	return space, true, nil
}

// single page
func FindPage(db *sql.DB, slug string) (domain.Page, bool, error) {
	var page domain.Page

	err := db.QueryRow(
		`
        SELECT title, slug, space_slug, content
        FROM pages
        WHERE slug = $1
        `,
		slug,
	).Scan(
		&page.Title,
		&page.Slug,
		&page.SpaceSlug,
		&page.Content,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Page{}, false, nil
		}

		return domain.Page{}, false, err
	}

	return page, true, nil
}

// page by space
func FindPagesBySpace(db *sql.DB, spaceSlug string) ([]domain.Page, error) {
	rows, err := db.Query(
		`
		SELECT title, slug, space_slug, content
		FROM pages
		WHERE space_slug = $1
		ORDER BY title
		`,
		spaceSlug,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []domain.Page

	for rows.Next() {
		var page domain.Page

		err := rows.Scan(
			&page.Title,
			&page.Slug,
			&page.SpaceSlug,
			&page.Content,
		)
		if err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pages, nil
}

// spaces

func FindSpaces(db *sql.DB) ([]domain.Space, error) {
	rows, err := db.Query(
		`
		SELECT name, slug
		FROM spaces
		ORDER BY name
		
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spaces []domain.Space

	for rows.Next() {
		var space domain.Space

		err := rows.Scan(
			&space.Name,
			&space.Slug,
		)
		if err != nil {
			return nil, err
		}
		spaces = append(spaces, space)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return spaces, nil
}
