# Pokedex

Pokedex is a project based on the [Boot.dev Pokedex CLI course](https://www.boot.dev/courses/build-pokedex-cli-golang), leveraging the [PokeAPI](https://pokeapi.co/) to provide Pokémon data. While I'm not a Pokémon fan, this project has been an engaging way to build skills in project structure, package management, dependency injection, and caching in Go.

## What did I learn?

I started my Go journey about a year ago and try to put in some work each day on a project. While working through the [Boot.dev](https://www.boot.dev) courses, this one challenged me a little bit more. The more challenging aspects were around implemeting the caching. I spent a bit more time on my own structuring the project in a way that made a little more sense.

## Features

- Explore the map and different areas
- Find Pokémon in specific areas
- Catch Pokémon with a simple probability formula
- List all caught Pokémon
- Inspect stats for caught Pokémon
- Command history support (up arrow)
- Persistent storage for caught Pokémon

## Available Commands

```bash
map                # Show the entire map by pulling from the API; typing 'map' again advances the map forward (pagination)
mapb               # Show the previous area
explore 'area'     # Show which Pokémon are in a specific area
inspect 'pokemon'  # Show stats for a specific Pokémon (if it has been caught)
catch 'pokemon'    # Attempt to catch a specific Pokémon (currently uses a simple 25% chance formula)
exit               # Exit the program
help               # Display help for the program
pokedex            # Display all caught Pokémon
```

## Improvements

While I enjoyed the course and the challenge, I wanted to improve several aspects before considering the project complete. The following areas have been addressed:

- [ ] Test Coverage: Ensure all key components have test coverage
- [x] Persistence: Caught Pokémon persist between application runs
- [x] Code Refactor / Organization: The base project was reorganized into packages for API, cache, and REPL layers
- [x] Command History: Support for recalling previous commands (up arrow; note: 'history' command implemented, up arrow may require external dependencies)
- [ ] Pokémon Battles: Allow for battling between caught Pokémon
- [x] Help Output: Align 'help' output cleanly in the terminal

## Installation

```bash
# Clone the repository
git clone https://github.com/joshhartwig/pokedex.git

# Navigate to the project directory
cd pokedex
```

## Technologies Used

- Go

## Acknowledgments

- Pokémon data provided by [PokeAPI](https://pokeapi.co/)
- Inspired by the [Boot.dev Pokedex CLI course](https://www.boot.dev/courses/build-pokedex-cli-golang)
