-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    song VARCHAR(250) NOT NULL,
    group_name VARCHAR(250) NOT NULL,
    text TEXT,
    link VARCHAR(128),
    date_release VARCHAR(16),
    UNIQUE(song, group_name)
);
CREATE INDEX songs_song ON songs(song);
CREATE INDEX songs_group ON songs(group_name);
CREATE INDEX songs_date ON songs(date_release);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX songs_date;
DROP INDEX songs_group;
DROP INDEX songs_song;
DROP TABLE IF EXISTS songs;
-- +goose StatementEnd
