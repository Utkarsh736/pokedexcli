# Pokedex CLI (Go)
A small, interactive Pokédex command-line app built in Go that lets you explore the Pokémon world, discover encounters in location areas, catch Pokémon with a probability system, and manage/view your personal Pokédex — all from a simple REPL.

This project uses the public PokéAPI (v2) for Pokémon and location-area data and includes an in-memory cache to avoid repeated network calls and keep the CLI snappy.
​

## Demo (example session)
```text
Pokedex > help
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
map: Display next 20 location areas
mapb: Display previous 20 location areas
explore: Explore a location area
catch: Attempt to catch a pokemon
inspect: View details about a caught Pokemon
pokedex: View all caught Pokemon

Pokedex > map
... (20 location areas)

Pokedex > explore pastoria-city-area
Exploring pastoria-city-area...
Found Pokemon:
 - tentacool
 - tentacruel
 ...

Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu was caught!

Pokedex > inspect pikachu
Name: pikachu
Height: 4
Weight: 60
Stats:
  -hp: 35
  -attack: 55
  ...
Types:
  - electric

Pokedex > pokedex
Your Pokedex:
  - pikachu
```
## Commands
- `help`: Show all available commands and short descriptions.
- `exit`: Exit the CLI.
- `map`: Show the next page (20) of location areas (pagination).
- `mapb`: Show the previous page (20) of location areas (pagination).
- `explore <location_area>`: List Pokémon that can be encountered in the given location area. (Uses the PokéAPI “location-area” resource.)​
- `catch <pokemon>`: Attempt to catch a Pokémon; catch chance scales with difficulty using Pokémon data from PokéAPI (notably base experience).
- `inspect <pokemon>`: Show details of a Pokémon only if it has been caught (no API call needed).
- `pokedex`: List all caught Pokémon in your collection.

## Installation & Run
### Prerequisites
- Go installed (Go modules enabled).

### Run locally
```bash
git clone <your-repo-url>
cd pokedexcli
go test ./...
go run .
```
### Build a binary
```bash
go build
./pokedexcli
```
## Project structure
Typical layout (names may vary slightly depending on how you organized files):
- `main.go`: Program entry point; wires up the REPL loop and command dispatch.
- `repl.go` (or similar): Input cleaning/parsing (`cleanInput`) and REPL helpers.
- `commands.go` (and/or `command_*.go`): Command registry + command callbacks.
- `internal/pokeapi/`: Small client layer for PokéAPI calls and JSON decoding.
​- 'internal/pokecache/`: In-memory cache with TTL-like reaping using a ticker + mutex for concurrent safety.

## How it works (high-level)
- <b>REPL loop</b>: Reads a line from stdin, normalizes it (trim/lowercase/split), then dispatches by the first token (command).
- <b>Command registry</b>: Commands are registered in a map (name → metadata + callback), so adding new commands is clean and centralized.
- <b>Pagination state</b>: map / mapb maintain next and previous URLs so you can page through location areas.
- <b>Caching</b>: API responses are cached by URL key; re-visiting pages/areas is fast and reduces repeated API calls (aligned with PokéAPI fair-use guidance).
- <b>Pokédex state</b>: Caught Pokémon are stored in-memory (e.g., map[string]Pokemon) and used by inspect/pokedex without additional network requests.

## Testing
Run all tests from repo root:
```bash
go test ./...
```
What’s covered:
- Unit tests for input cleaning (cleanInput).
- Unit tests for cache behavior (Add/Get and TTL reaping).

## Skills demonstrated
- Building an interactive CLI (REPL) in Go
- Table-driven unit testing
- HTTP GET + JSON decoding against a real public API (PokéAPI v2)
​- Pagination using next / previous URLs
​- In-memory caching with:
  - mutex-protected map (concurrency safety)
  - background reaper using a ticker
- Practical Go project organization with internal packages

## Contributing
Contributions are welcome, especially around:
- More tests (commands, cache edge cases, input parsing)
- Better UX (command history, nicer formatting, autocomplete)
- Refactors that improve readability/testability

Suggested workflow:
1. Fork the repo
2. Create a feature branch
3. Add/adjust tests where appropriate
4. Open a PR with a short description and rationale

## Ideas for extensions
If you want to turn this into a portfolio-worthy project:

- Persistent storage (save/load Pokédex between runs)
- Command history (up arrow cycling)
- Random encounters and/or battle simulation
- Party system + leveling/evolution timers
- Different Poké Ball types with different catch modifiers
- More exploration UX (choices instead of typing long area names)
- More unit/integration tests and better dependency injection for API/client layers

## Notes / Credits
Data provided by PokéAPI v2 (public, no-auth, GET-only). Please be considerate and cache responses where possible.
