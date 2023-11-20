package juicer

type RayDirection uint8

const (
	Nort RayDirection = iota // north
	NoEa                     // north-east
	East                     // east
	SoEa                     // south-east
	Sout                     // south
	SoWe                     // south-west
	West                     // west
	NoWe                     // north-west
)

type DirectionIncrement int8

const (
	IncrementNort DirectionIncrement = 8  // north
	IncrementNoEa DirectionIncrement = 9  // north-east
	IncrementEast DirectionIncrement = 1  // east
	IncrementSoEa DirectionIncrement = -1 // south-east
	IncrementSout DirectionIncrement = -8 // south
	IncrementSoWe DirectionIncrement = -9 // south-west
	IncrementWest DirectionIncrement = -1 // west
	IncrementNoWe DirectionIncrement = 7  // north-west
)
