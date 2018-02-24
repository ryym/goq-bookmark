CREATE TABLE users (
    id SERIAL NOT NULL,
    name VARCHAR(128) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    UNIQUE (name)
);

CREATE TABLE entries (
    id SERIAL NOT NULL,
    url VARCHAR(512) NOT NULL,
    title VARCHAR(256) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    UNIQUE (url)
);

CREATE TABLE bookmarks (
    id SERIAL NOT NULL,
    user_id INT NOT NULL,
    entry_id INT NOT NULL,
    comment VARCHAR(256) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    UNIQUE (user_id, entry_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (entry_id) REFERENCES entries (id)
);

