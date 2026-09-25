package domain

func SampleSpaces() []Space {
	return []Space{
		{
			Name: "Infrastructure",
			Slug: "infrastructure",
		},
		{
			Name: "Django",
			Slug: "django",
		},
		{
			Name: "Databases",
			Slug: "databases",
		},
		{
			Name: "Security",
			Slug: "security",
		},
	}
}

func SampePages() []Page {
	return []Page{
		{
			Title:     "PostgreSQL Migration",
			Slug:      "postgresql-migration",
			SpaceSlug: "databases",
			Content:   "This guide covers the basic process for migrating a PostgreSQL database safely.",
		},
		{
			Title:     "Backup and Restore",
			Slug:      "backup-restore",
			SpaceSlug: "databases",
			Content:   "This guide covers database backup and restore procedures.",
		},
		{
			Title:     "Docker Deployment",
			Slug:      "docker-deployment",
			SpaceSlug: "infrastructure",
			Content:   "This guide covers the basic process for deploying an application with Docker.",
		},
		{
			Title:     "Django Deployment",
			Slug:      "django-deployment",
			SpaceSlug: "django",
			Content:   "This guide covers the basic steps for deploying a Django application.",
		},
		{
			Title:     "What’s new in Django 6.1",
			Slug:      "django-upgrade",
			SpaceSlug: "django",
			Content:   "These release notes cover the new features, as well as some backwards incompatible changes you’ll want to be aware of when upgrading from Django 6.0 or earlier",
		},
	}
}

func FindSpace(slug string) (Space, bool) {
	for _, space := range SampleSpaces() {
		if space.Slug == slug {
			return space, true
		}
	}
	return Space{}, false
}

func FindPagesBySpace(spaceSlug string) []Page {
	var pages []Page

	for _, page := range SampePages() {
		if page.SpaceSlug == spaceSlug {
			pages = append(pages, page)
		}
	}
	return pages
}

func FindPage(slug string) (Page, bool) {
	for _, page := range SampePages() {
		if page.Slug == slug {
			return page, true
		}
	}
	return Page{}, false
}
