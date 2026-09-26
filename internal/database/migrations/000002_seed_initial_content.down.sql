DELETE FROM pages
WHERE slug IN (
    'postgresql-migration',
    'backup-restore',
    'docker-deployment',
    'django-deployment'
);

DELETE FROM spaces
WHERE slug IN (
    'infrastructure',
    'django',
    'databases',
    'security'
);