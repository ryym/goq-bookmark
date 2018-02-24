INSERT INTO users (id, name) VALUES
(1, 'Alice'),
(2, 'Bob')
;
SELECT setval('users_id_seq', 2);

INSERT INTO entries (id, url, title) VALUES
(1, 'http://example.com', 'The example domain'),
(2, 'https://github.com', 'GitHub'),
(3, 'https://gist.github.com', 'Gist'),
(4, 'https://google.co.jp', 'Google'),
(5, 'https://amazon.co.jp', 'Amazon')
;
SELECT setval('entries_id_seq', 5);

INSERT INTO bookmarks (id, user_id, entry_id, comment) VALUES
(1, 1, 2, 'So many great projects!'),
(2, 1, 3, ''),
(3, 1, 5, 'I have a prime membership'),
(4, 2, 1, 'The example domain')
;
SELECT setval('bookmarks_id_seq', 4);
