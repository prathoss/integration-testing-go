CREATE TABLE pictures (
    id SERIAL PRIMARY KEY,
    url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    author_id INTEGER NOT NULL,
    FOREIGN KEY (author_id) REFERENCES profiles(id)
);

CREATE INDEX idx_pictures_author ON pictures(author_id);
