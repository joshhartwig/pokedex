# Pokedex

Pokedex is a project from bootdev[https://www.boot.dev/courses/build-pokedex-cli-golang] leveraging the pokedex[https://pokeapi.co/] api. I don't know or care about Pokemon but the project is engaging and builds some skills about project, package management, dependency injection and caching.

## Improvements

The base project has the following commands

```bash
map # Show the entire map pulling from the api, typing in map again advances the map forward paging the results
mapb # Show the previous area
explore 'area' # Show which pokemon are in a specific area
inspect 'pokemon' # Show stats for a specific pokemon (assuming it is caught)
catch 'pokemon' # Attempt to catch a specific pokemon. For now uses simple 25% formula
exit # Exit the program 
help # Display help for the program
pokedex # Displays all caught pokemon
-
```

## Improvements

While I enjoyed the course and the challenge I wanted to improve on quite a few things before calling the project a wrap. The following areas I improved:

[ ] Test Coverage: Ensure we have test coverage for all key components
[ ] Persistence: Caught Pokemon persist from app run to app run
[x] Code refactor / organization: The base project was a bit of a jumbled mess. Organize the api, cache layer and repl into packages
[ ] Support for up arrow: Allow for previous command to be remembered
[ ] Allow for battling between caught Pokemon
[x] Align 'help' output cleanly in terminal

## Features

- Pokemon search functionality
- Detailed Pokemon information
- Type matchups
- Evolution chains

## Installation

```bash
# Clone the repository
git clone https://github.com/joshhartwig/pokedex/pokedex.git

# Navigate to the project directory
cd pokedex


## Technologies Used

- Go!

## Acknowledgments

- Pokemon data provided by [PokeAPI](https://pokeapi.co/)
