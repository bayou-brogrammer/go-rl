# Go-RL Code Guidelines

## Build & Run Commands
- `make build` - Build application to build/game
- `make run` - Build and run
- `make run-direct` - Run without building
- `make test` - Run all tests
- `go test ./game/...` - Run specific tests
- `make fmt` - Format code
- `make vet` - Lint code
- `make deps` - Install dependencies
- `make clean` - Clean build artifacts

## Code Style Guidelines

### Imports
- Group standard library imports first, then third-party packages, finally local packages
- Sort alphabetically within groups

### Error Handling
- Use the `logerror` package for user-facing errors
- Ensure error messages are capitalized
- Return errors rather than handling them prematurely

### Naming & Types
- Use camelCase for variable names, PascalCase for exported functions/types
- Define interfaces for behaviors (e.g., Consumable, Targetter)
- Use meaningful type names that describe purpose (e.g., HealingPotion)

### Documentation
- Document public functions and types with comments
- Include usage examples for complex APIs