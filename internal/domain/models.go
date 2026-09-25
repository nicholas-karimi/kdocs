package domain

// Space/Page
// represent actual things that exist in KDocs.

type Space struct {
	Name string
	Slug string
}

type Page struct {
	Title     string
	Slug      string
	SpaceSlug string
	Content   string
}
