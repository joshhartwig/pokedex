-- name: AddPokemon :one
INSERT INTO pokemon (id, pokemon_name, json_data)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPokemonByName :one
SELECT * FROM pokemon
WHERE pokemon_name = $1
LIMIT 1;

-- name: ListPokemon :many
SELECT * FROM pokemon
ORDER BY caught_at DESC;