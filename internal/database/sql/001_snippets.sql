-- +goose Up
CREATE TABLE IF NOT EXISTS snippets (
    id INTEGER PRIMARY KEY,
    name TEXT,
    text TEXT
);

INSERT INTO snippets (name, text) VALUES
    ('Delete Git Branches', 'git branch | grep -v \"^\\*\" | xargs git branch -d'),
    ('Remove Git Tracking', 'rm -rf .git');

-- +goose Down
DROP TABLE snippets;
