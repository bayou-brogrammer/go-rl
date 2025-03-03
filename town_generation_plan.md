# Town Generation Method Plan

## Overview

The town generation method will create a town with the following characteristics:

- Grid-based streets forming the main structure
- Some organic elements (like slightly curved streets or irregular building placements)
- Walled districts dividing the town into sections
- A central town square or marketplace

## Implementation Plan

### 1. File Structure

We'll create a new file called `town.go` in the `game/dungeon/` directory that will contain our town generation code.

```plaintext
game/dungeon/
├── bsp.go
├── dungeon.go
└── town.go (new)
```

### 2. Town Generation Algorithm

The algorithm will work in several phases:

1. Initialize Map with Walls
2. Create Town Layout
3. Generate Streets
4. Create Districts
5. Add Town Walls
6. Place Buildings
7. Create Town Square/Marketplace
8. Add Entrances/Exits
9. Ensure Connectivity

### 3. Detailed Implementation Steps

#### 3.1. Town Structure Definition

We'll define structures to represent town elements:

```go
// TownDistrict represents a district in the town
type TownDistrict struct {
    X, Y          int // Top-left corner
    Width, Height int
    Type          DistrictType
}

// DistrictType represents different types of districts
type DistrictType int

const (
    ResidentialDistrict DistrictType = iota
    CommercialDistrict
    NobleDistrict
)

// TownBuilding represents a building in the town
type TownBuilding struct {
    X, Y          int // Top-left corner
    Width, Height int
    Type          BuildingType
}

// BuildingType represents different types of buildings
type BuildingType int

const (
    House BuildingType = iota
    Shop
    Inn
    Temple
)
```

#### 3.2. Main Generation Function

```go
// GenerateTown generates a town-like map
func (m *Map) GenerateTown() {
    // Initialize map with walls
    size := m.Grid.Size()
    for x := range size.X {
        for y := range size.Y {
            m.Grid.Set(gruid.Point{X: x, Y: y}, Wall)
        }
    }

    // Generate town layout
    districts := generateDistricts(m)

    // Generate streets
    streets := generateStreets(m, districts)

    // Place buildings
    buildings := placeBuildings(m, districts, streets)

    // Create town square/marketplace
    createTownSquare(m, districts)

    // Add town walls
    addTownWalls(m)

    // Ensure connectivity
    ensureConnectivity(m)
}
```

#### 3.3. District Generation

```go
// generateDistricts divides the map into districts
func generateDistricts(m *Map) []TownDistrict {
    // Divide the map into 2-4 districts
    // Each district will have its own characteristics
    // Return the list of districts
}
```

#### 3.4. Street Generation

```go
// generateStreets creates a network of streets
func generateStreets(m *Map, districts []TownDistrict) []gruid.Path {
    // Create main streets between districts
    // Create grid-based streets within districts
    // Add some organic variations to streets
    // Return the list of street paths
}
```

#### 3.5. Building Placement

```go
// placeBuildings adds buildings to the town
func placeBuildings(m *Map, districts []TownDistrict, streets []gruid.Path) []TownBuilding {
    // Place buildings along streets
    // Vary building sizes based on district type
    // Ensure buildings don't overlap with streets
    // Return the list of buildings
}
```

#### 3.6. Town Square Creation

```go
// createTownSquare adds a town square or marketplace
func createTownSquare(m *Map, districts []TownDistrict) {
    // Find a suitable location for the town square
    // Create an open area
    // Add some decorative elements
}
```

#### 3.7. Town Walls

```go
// addTownWalls surrounds the town with walls
func addTownWalls(m *Map) {
    // Create outer walls around the town
    // Add gates at street exits
}
```

#### 3.8. Connectivity Check

```go
// ensureConnectivity makes sure all floor tiles are connected
func ensureConnectivity(m *Map) {
    // Use existing connectivity check code from dungeon.go
    // Fill in disconnected areas or create paths to connect them
}
```

### 4. Integration with Existing Code

We'll modify the `Generate` method in `dungeon.go` to allow calling our town generation function:

```go
// Generate fills the Grid attribute of m with a procedurally generated map.
func (m *Map) Generate() {
    // Decide which generation method to use
    // For now, we'll just use the town generation
    m.GenerateTown()

    // Later, this could be expanded to choose between different generation methods
    // based on parameters or random selection
}
```

## Technical Considerations

### 1. Random Number Generation

We'll use the existing random number generator from the Map struct:

```go
rng := m.rand
```

### 2. Map Representation

We'll use the existing Grid structure from the Map struct:

```go
m.Grid.Set(gruid.Point{X: x, Y: y}, Floor)
```

### 3. Connectivity

We'll use the existing connectivity check code from dungeon.go to ensure all floor tiles are connected:

```go
freep := m.RandomFloor()
pr := paths.NewPathRange(m.Grid.Range())
pr.CCMap(&path{m: m}, freep)
```

### 4. Tile Types

We'll use the existing tile types (Wall, Floor) and potentially add new ones if needed:

```go
const (
    Wall rl.Cell = iota
    Floor
    Door  // New tile type for doors
)
```

## Visual Representation

Here's a conceptual representation of what the generated town might look like:

```plaintext
##############################################
#                                            #
#   ################    ###################  #
#   #              #    #                 #  #
#   #  DISTRICT 1  #    #   DISTRICT 2    #  #
#   #              #    #                 #  #
#   #  [][][][][]  #    #   [][][][][]    #  #
#   #  [][][][][]  #    #   [][][][][]    #  #
#   #  [][][][][]  #    #   [][][][][]    #  #
#   ################    ###################  #
#                                            #
#                 [][]                       #
#                 [][]                       #
#                 SQUARE                     #
#                                            #
#   ################    ###################  #
#   #              #    #                 #  #
#   #  DISTRICT 3  #    #   DISTRICT 4    #  #
#   #              #    #                 #  #
#   #  [][][][][]  #    #   [][][][][]    #  #
#   #  [][][][][]  #    #   [][][][][]    #  #
#   #  [][][][][]  #    #   [][][][][]    #  #
#   ################    ###################  #
#                                            #
##############################################
```

Where:

- `#` represents walls
- `[]` represents buildings
- Empty spaces represent streets
- `SQUARE` represents the town square/marketplace
- `DISTRICT X` represents different districts

## Testing Strategy

1. Visual inspection of generated towns
2. Connectivity testing to ensure all floor tiles are reachable
3. Performance testing to ensure generation is efficient
