CREATE TABLE picture_views (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL,
    picture_id INTEGER NOT NULL,
    view_count INTEGER NOT NULL DEFAULT 1,
    last_viewed_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY (profile_id) REFERENCES profiles(id),
    FOREIGN KEY (picture_id) REFERENCES pictures(id),
    UNIQUE (profile_id, picture_id)
);

CREATE INDEX idx_picture_views_profile ON picture_views(profile_id);
CREATE INDEX idx_picture_views_picture ON picture_views(picture_id);
