snippetbox=# INSERT INTO snippets (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.',
    current_timestamp,
    current_timestamp + INTERVAL '365 day'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– N',
    current_timestamp,
    current_timestamp + INTERVAL '365 day'
);
INSERT INTO snippets (title, content, created, expires) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\',
    current_timestamp,
    current_timestamp + INTERVAL '365 day'
);