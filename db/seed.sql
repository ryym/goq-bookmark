INSERT INTO users (id, name) VALUES
(1, 'Alice')
;
SELECT setval('users_id_seq', 1);

INSERT INTO entries (id, url, title) VALUES
(1, 'http://example.com', 'The example domain'),
(2, 'https://github.com', 'GitHub'),
(3, 'https://gist.github.com', 'Gist')
;
SELECT setval('entries_id_seq', 3);

INSERT INTO bookmarks (id, user_id, entry_id, comment) VALUES
(1, 1, 1, NULL),
(2, 1, 2, 'So many great projects!')
;
SELECT setval('bookmarks_id_seq', 2);
