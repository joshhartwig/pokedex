-- +goose Up

CREATE TABLE pokemon (
    id UUID PRIMARY KEY,
    pokemon_name TEXT NOT NULL,
    caught_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    json_data json,
    UNIQUE(pokemon_name)
);

-- +goose Down
DROP TABLE pokemon;