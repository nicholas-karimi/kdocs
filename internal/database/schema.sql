
CREATE TABLE spaces(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE 
);

CREATE TABLE pages (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    space_slug TEXT NOT NULL REFERENCES spaces(slug),
    content TEXT NOT NULL
);