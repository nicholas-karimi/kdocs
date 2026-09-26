INSERT INTO spaces (name, slug)
VALUES
    ('Infrastructure', 'infrastructure'),
    ('Django', 'django'),
    ('Databases', 'databases'),
    ('Security', 'security');

INSERT INTO pages (title, slug, space_slug, content)
VALUES
    (
        'PostgreSQL Migration',
        'postgresql-migration',
        'databases',
        'This guide covers the basic process for migrating a PostgreSQL database safely.'
    ),
    (
        'Backup and Restore',
        'backup-restore',
        'databases',
        'This guide covers database backup and restore procedures.'
    ),
    (
        'Docker Deployment',
        'docker-deployment',
        'infrastructure',
        'This guide covers the basic process for deploying an application with Docker.'
    ),
    (
        'Django Deployment',
        'django-deployment',
        'django',
        'This guide covers the basic steps for deploying a Django application.'
    );